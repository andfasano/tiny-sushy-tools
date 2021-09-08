package redfish

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/gorilla/mux"
)

//Server represents a mock server supporting partially the RedFish protocol
type Server struct {
	router *mux.Router

	systems map[string]*system

	TinySushyPort string
	TinyOobUser   string
	TinyOobIP     string
	TinyOobKey    string
}

//New creates a new instance of the Redfish server
func New() *Server {
	return &Server{
		systems: make(map[string]*system),
	}
}

//Check if valid Port
func isValidPort(port string) (bool, error) {
	_, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return false, err
	}
	return true, nil
}

//Sanitize username
func isValidUsername(user string) (bool, error) {
	//NOT IMPLEMENTED - Do we have to sanitize username?
	return true, nil
}

//Check if valid IP
func isValidIP(ip string) (bool, error) {
	if net.ParseIP(ip) == nil {
		return false, nil
	}
	return true, nil
}

//Check if file exist
func isValidKeyPath(keypath string) (bool, error) {
	//NOTE: This does NOT work for check of relative path
	// _, err := os.Stat(keypath)
	// if err != nil {
	// 	return false, err
	// }
	return true, nil
}

//Start initialize and kicks off the redfish server
func (rf *Server) Start() {
	rf.router = mux.NewRouter()

	//Run validations for variables
	if isValid, _ := isValidPort(rf.TinySushyPort); !isValid {
		log.Println("Invalid port ", rf.TinySushyPort)
		log.Fatal("Invalid port")
	}
	if isValid, _ := isValidUsername(rf.TinyOobUser); !isValid {
		log.Println("Invalid username string ", rf.TinyOobUser)
		log.Fatal("Invalid username string")
	}
	if isValid, _ := isValidIP(rf.TinyOobIP); !isValid {
		log.Println("Invalid IP address ", rf.TinyOobIP)
		log.Fatal("Invalid IP address")
	}
	if isValid, _ := isValidKeyPath(rf.TinyOobKey); !isValid {
		log.Println("File not found ", rf.TinyOobKey)
		log.Fatal("File not found (private key)")
	}

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

	log.Println("Starting RedFish mock server on port ", rf.TinySushyPort)
	log.Fatal(http.ListenAndServe(":"+rf.TinySushyPort, rf.router))
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

				errorTemplate := `
				{
					"error": {
					  "code": "Base.1.0.GeneralError",
					  "message": "A general error has occurred. See ExtendedInfo for more information.",
					  "@Message.ExtendedInfo": [
						{
						  "MessageId": "GEN1234",
						  "RelatedProperties": [],
						  "Message": "Unable to process the request because an error occurred.",
						  "MessageArgs": [],
						  "Severity": "Critical",
						  "Resolution": "Retry the operation. If the issue persists, contact your system administrator."
						}
					  ]
					}
				  }`

				http.Error(w, errorTemplate, http.StatusUnauthorized)
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

	systemID := mux.Vars(r)["identity"]

	switch r.Method {
	case http.MethodGet:

		s, ok := rf.systems[systemID]
		if ok == false {
			s = newSystem(systemID, rf)
			rf.systems[systemID] = s
		}

		if !rf.checkBMCCredentials(systemID, w, r) {
			return
		}
		s.Send(w)

		break
	default:
		log.Fatal("Method not supported")
	}
}

func (rf *Server) handleSystems(w http.ResponseWriter, r *http.Request) {
	rf.logRequest("/redfish/v1/Systems", r)
	response := "Listing not allowed. Please specify system identity /redfish/v1/Systems/{identity}"
	http.Error(w, response, http.StatusForbidden)
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
