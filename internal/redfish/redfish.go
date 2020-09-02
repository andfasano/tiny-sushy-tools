package redfish

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

//Server represents a mock server supporting partially the RedFish protocol
type Server struct {
	router *mux.Router

	systems map[string]*system
}

//New creates a new instance of the Redfish server
func New() *Server {
	return &Server{
		systems: make(map[string]*system),
	}
}

//Start initialize and kicks off the redfish server
func (rf *Server) Start(port string) {
	rf.router = mux.NewRouter()

	//RedFish protocol
	rf.router.HandleFunc("/", rf.handleCatchAll)
	rf.router.HandleFunc("/redfish/v1/", rf.handleEntrypoint)

	// r.router.HandleFunc("/redfish/v1/Chassis", handleChassis)
	// r.router.HandleFunc("/redfish/v1/Chassis/{identity}", handleChassisByID).Methods("GET", "PATCH")
	// r.router.HandleFunc("/redfish/v1/Chassis/{identity}/Thermal", handleChassisByIDThermal).Methods("GET")

	// r.router.HandleFunc("/redfish/v1/Managers", handleManagers)
	// r.router.HandleFunc("/redfish/v1/Managers/{identity}", handleManagersByID).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Managers/{identity}/VirtualMedia", handleManagersByIDVirtualMedia).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Managers/{identity}/VirtualMedia/{device}", handleManagersByIDVirtualMediaDevice).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Managers/{identity}/VirtualMedia/{device}/Actions/VirtualMedia.InsertMedia", handleVirtualMediaActionInsertMedia).Methods("POST")
	// r.router.HandleFunc("/redfish/v1/Managers/{identity}/VirtualMedia/{device}/Actions/VirtualMedia.EjectMedia", handleVirtualMediaActionEjectMedia).Methods("POST")

	rf.router.HandleFunc("/redfish/v1/Systems", rf.handleSystems)
	rf.router.HandleFunc("/redfish/v1/Systems/{identity}", rf.handleSystemsByID).Methods("GET", "PATCH")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/EthernetInterfaces", handleSystemsEthernetInterfaces).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/EthernetInterfaces/{nic_id}", handleSystemsEthernetInterfacesByNicID).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/Actions/ComputerSystem.Reset", handleSystemsActionReset).Methods("POST")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/BIOS", handleSystemsBIOS).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/BIOS/Settings", handleSystemsBIOSSettings).Methods("GET", "PATCH")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/BIOS/Actions/Bios.ResetBios", handleSystemsActionResetBIOS).Methods("POST")

	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/SimpleStorage", handleSystemsSimpleStorage).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/SimpleStorage/{simple_storage_id}", handleSystemsSimpleStorageByID).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/Storage", handleSystemsStorage).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/Storage/{storage_id}", handleSystemsStorageByID).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/Storage/{storage_id}/Drives/{drive_id}", handleSystemsStorageDrivesByID).Methods("GET")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/Storage/{storage_id}/Volumes", handleSystemsStorageVolumes).Methods("GET", "POST")
	// r.router.HandleFunc("/redfish/v1/Systems/{identity}/Storage/{storage_id}/Volumes/{volume_id}", handleSystemsStorageVolumesByID).Methods("GET")

	//Mock protocol
	rf.router.HandleFunc("/mock/Systems/{identity}/Credentials", rf.handleMockSystemsCredentials).Methods("PUT")

	log.Println("Starting RedFish mock server on port ", port)
	log.Fatal(http.ListenAndServe(":"+port, rf.router))
}

func (rf *Server) logRequest(src string, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(src, " --- ", string(requestDump))
}

func (rf *Server) checkBMCCredentials(UUID string, w http.ResponseWriter, r *http.Request) (match bool) {
	if username, password, ok := r.BasicAuth(); ok {
		if s, ok := rf.systems[UUID]; ok {
			match = s.Username == username && s.Password == password
			if !match {
				http.Error(w, "Credentials mismatch", http.StatusBadRequest)
			}
		} else {
			http.Error(w, "Unable to find system "+UUID, http.StatusBadRequest)
		}
	} else {
		log.Println("No auth")
		match = true
	}

	return
}

func (rf *Server) handleMockSystemsCredentials(w http.ResponseWriter, r *http.Request) {
	UUID := mux.Vars(r)["identity"]
	log.Println("### Changing system credentials for " + UUID)

	creds := struct {
		Username string
		Password string
	}{}

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		http.Error(w, "Missing one or more parameters", http.StatusBadRequest)
		return
	}

	if s, ok := rf.systems[UUID]; ok {
		s.Username = creds.Username
		s.Password = creds.Password
		fmt.Fprintf(w, "Creds: %+v", rf.systems[UUID])
	} else {
		http.Error(w, "Unable to find system "+UUID, http.StatusBadRequest)
		return
	}

}

func (rf *Server) handleSystemsByID(w http.ResponseWriter, r *http.Request) {
	log.Println("-- Request System " + mux.Vars(r)["identity"])

	UUID := mux.Vars(r)["identity"]
	if !rf.checkBMCCredentials(UUID, w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:

		s, ok := rf.systems[UUID]
		if ok == false {
			s = newSystem(UUID)
			rf.systems[UUID] = s
		}
		s.Send(w)

		break
	default:
		log.Fatal("Method not supported")
	}
}

func (rf *Server) handleSystems(w http.ResponseWriter, r *http.Request) {
	rf.logRequest("/redfish/v1/Systems", r)
}

func (rf *Server) handleEntrypoint(w http.ResponseWriter, r *http.Request) {
	log.Println("-- Main entry")

	rootTemplate := `{
	"@odata.type": "#ServiceRoot.v1_0_2.ServiceRoot",
	"Id": "RedvirtService",
	"Name": "Redvirt Service",
	"RedfishVersion": "1.0.2",
	"UUID": "85775665-c110-4b85-8989-e6162170b3ec",
	"Systems": {
		"@odata.id": "/redfish/v1/Systems"
	},
	"Managers": {
		"@odata.id": "/redfish/v1/Managers"
	},
	"@odata.id": "/redfish/v1/",
	"@Redfish.Copyright": "Copyright 2014-2016 Distributed Management Task Force, Inc. (DMTF). For the full DMTF copyright policy, see http://www.dmtf.org/about/policies/copyright."
}`

	w.Write([]byte(rootTemplate))
}

func (rf *Server) handleCatchAll(w http.ResponseWriter, r *http.Request) {
	rf.logRequest("/", r)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
