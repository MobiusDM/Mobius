package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/notawar/mobius/mobius-server/api"
	"github.com/notawar/mobius/mobius-server/pkg/service"
	"github.com/notawar/mobius/mobius-server/pkg/websocket"
)

func main() {
	fmt.Println("🚀 Starting Mobius MDM Test Server...")

	// Create WebSocket hub for real-time communication
	wsHub := websocket.NewHub()
	ctx := context.Background()
	go wsHub.Run(ctx)

	// Create service implementations
	deviceService := service.NewDeviceService()
	policyService := service.NewPolicyService()
	applicationService := service.NewApplicationService()
	authService := service.NewAuthService()
	groupService := service.NewGroupService()

	// Create API dependencies with WebSocket support
	deps := &api.Dependencies{
		DeviceService:      deviceService,
		PolicyService:      policyService,
		GroupService:       groupService,
		ApplicationService: applicationService,
		AuthService:        authService,
		WSHub:             wsHub,
	}

	// Create router
	router := api.NewRouter(deps)

	// Create server
	server := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("🌐 Server starting on http://localhost:8081\n")
		fmt.Printf("📋 Health check: http://localhost:8081/api/v1/health\n")
		fmt.Printf("📖 API docs: http://localhost:8081/api/v1/\n")
		fmt.Printf("🔌 WebSocket endpoint: ws://localhost:8081/api/v1/ws\n")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server exited cleanly")
}
