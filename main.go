package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mattjw79/aws-api/internal/amazonws"
)

type API struct {
}

func main() {
	api := API{}
	router := mux.NewRouter()
	router.
		Methods("GET").
		Path("/api/v1/regions").
		HandlerFunc(api.regions)
	router.
		Methods("GET").
		Path("/api/v1/region/{region:[^/?]*}").
		HandlerFunc(api.regions)
	router.
		Methods("GET").
		Path("/api/v1/instances").
		HandlerFunc(api.instances)
	router.
		Methods("GET").
		Path("/api/v1/instance/{id:[^/?]*}").
		HandlerFunc(api.instanceByID)

	fmt.Println("server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func (a API) regions(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	regions, err := amazonws.Regions()
	if err != nil {
		fmt.Println(err)
	}

	if region == "" {
		json.NewEncoder(w).Encode(regions)
		return
	}

	for _, r := range regions {
		if region == r.Name {
			json.NewEncoder(w).Encode(r)
			return
		}
	}
}

func (a API) instanceByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	instances, err := amazonws.FindInstanceByID(id)
	if err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(
		struct {
			Count     int                 `json:"count"`
			Instances []amazonws.Instance `json:"instances"`
		}{
			Count:     len(instances),
			Instances: instances,
		},
	)

}

func (a API) instances(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	var regions []amazonws.Region
	if len(r.Form["region"]) == 0 {
		regions, err = amazonws.Regions()
		if err != nil {
			log.Printf("error reading regions: %s\n", err)
		}
	} else {
		for _, regionName := range r.Form["region"] {
			regions = append(regions, amazonws.NewRegion(regionName))
		}
	}

	names := r.Form["name"]

	var wgPub, wgSub sync.WaitGroup
	instanceChannel := make(chan amazonws.Instance, 50)

	for _, region := range regions {
		wgPub.Add(1)
		go func(r amazonws.Region) {
			defer wgPub.Done()
			var instances []amazonws.Instance
			if len(names) == 0 {
				instances, err = r.Instances()
				if err != nil {
					log.Printf("error reading instances from region %s: %s\n", r.Name, err)
					return
				}
			} else {
				for _, name := range names {
					i, _ := r.InstanceByName(name)
					instances = append(instances, i...)
				}
			}
			for _, instance := range instances {
				instanceChannel <- instance
			}
		}(region)
	}

	wgSub.Add(1)
	instances := make([]amazonws.Instance, 0, 10)
	go func() {
		defer wgSub.Done()
		for instance := range instanceChannel {
			instances = append(instances, instance)
		}
	}()

	wgPub.Wait()
	close(instanceChannel)

	wgSub.Wait()
	json.NewEncoder(w).Encode(
		struct {
			Count     int                 `json:"count"`
			Instances []amazonws.Instance `json:"instances"`
		}{
			Count:     len(instances),
			Instances: instances,
		},
	)
}
