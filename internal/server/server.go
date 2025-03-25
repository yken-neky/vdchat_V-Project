package server

import (
	"flag"
	"os"
	"time"

	"v/internal/handlers"
	w "v/pkg/webrtc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/gofiber/websocket/v2"
)

var (
	addr = flag.String("addr", ":"+os.Getenv("PORT"), "")
	cert = flag.String("cert", "", "")
	key  = flag.String("key", "", "")
)

func Run() error {
	flag.Parse()

	if *addr == ":" {
		*addr = ":8080"
	}

	engine := html.New("./views", ".html")
	engine.Reload(true) // Recargar templates en desarrollo

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main", // Asegurar que coincide con el nombre de la definición
	})

	app.Use(func(c *fiber.Ctx) error {
		if c.Protocol() == "http" {
			return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), 301)
		}
		return c.Next()
	})

	// Configurar Middleware de Seguridad Adicional
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://localhost, https://192.168.1.100",
		AllowMethods:     "GET,POST,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net ws:wss; "+
			"img-src 'self' data:; font-src 'self' https://cdn.jsdelivr.net")
		return c.Next()
	})
	// ------------------------------------------------------------

	app.Use(logger.New())
	app.Get("/", handlers.Welcome)
	app.Get("/room/create", handlers.RoomCreate)
	app.Get("/room/:uuid", handlers.Room)
	app.Get("/room/:uuid/websocket", websocket.New(handlers.RoomWebsocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))
	app.Get("/room/:uuid/chat", handlers.RoomChat)
	app.Get("/room/:uuid/chat/websocket", websocket.New(handlers.RoomChatWebsocket))
	app.Get("/room/:uuid/viewer/websocket", websocket.New(handlers.RoomViewerWebsocket))

	// Streams
	app.Get("/stream/:suuid", handlers.Stream)
	app.Get("/stream/:suuid/websocket", websocket.New(handlers.StreamWebsocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))
	app.Get("/stream/:suuid/chat/websocket", websocket.New(handlers.StreamChatWebsocket))
	app.Get("/stream/:suuid/viewer/websocket", websocket.New(handlers.StreamViewerWebsocket))

	app.Static("/stylesheets", "./assets/stylesheets")
	app.Static("/javascript", "./assets/javascript")
	app.Static("/", "./assets")

	w.Rooms = make(map[string]*w.Room)
	w.Streams = make(map[string]*w.Room)
	go dispatchKeyFrames()

	// Modificar el bloque final de Listen para mejor manejo de HTTPS
	if *cert != "" && *key != "" {
		return app.ListenTLS(":443", *cert, *key)
	}
	return app.Listen(*addr)
}

func dispatchKeyFrames() {
	for range time.NewTicker(time.Second * 3).C {
		for _, room := range w.Rooms {
			room.Peers.DispatchKeyFrame()
		}
	}
}
