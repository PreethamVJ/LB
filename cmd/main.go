package main

import (
	"log"

	"github.com/PreethamVJ/LB/internal/config"
	"github.com/PreethamVJ/LB/internal/loadbalancer"
)

// func main() {
// 	cfg, err := config.LoadConfig("test/config.toml")
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}
// 	fmt.Printf("%#v\n", cfg)
// 	fmt.Printf("Load Balancer Address: %s\n", cfg.LoadBalancer.Address) // Replace 'Address' with the actual field name, e.g., 'Host'
// 	fmt.Printf("Load Balancer Port: %d\n", cfg.LoadBalancer.Port)
// 	fmt.Printf("Load Balancer Algorithm: %s\n", cfg.LoadBalancer.Algorithm)

// 	fmt.Println("Servers:")
// 	for _, srv := range cfg.LoadBalancer.Servers {
// 		fmt.Printf(" - Address: %s, Port: %d, Max Connections: %d, Weight: %d\n",
// 			srv.Address, srv.Port, srv.MaxConnections, srv.Weight)
// 	}
// }

func main() {
	cfg, err := config.LoadConfig("test/config.toml")
	if err != nil {
		log.Fatal(err)
	}

	var lb loadbalancer.LoadBalancer
	switch cfg.LoadBalancer.Layer { // Add 'layer' field to config.toml
	case 4:
		lb = loadbalancer.NewLayer4(cfg)
	case 7:
		lb = loadbalancer.NewLayer7(cfg)
	default:
		log.Fatal("Unsupported layer")
	}
	log.Println("Starting load balancer...")
	if err := lb.Start(); err != nil {
		log.Fatal(err)
	}
}
