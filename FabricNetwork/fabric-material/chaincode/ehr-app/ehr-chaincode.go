
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Tuna struct {
	Vessel string `json:"vessel"`
	Timestamp string `json:"timestamp"`
	Location  string `json:"location"`
	Holder  string `json:"holder"`
}

type Medecin struct {
	NomPrenom string `json:"nomprenom"`
	Specialite  string `json:"specialite"`
	Cin  string `json:"cin"`
}

type Patient struct {
	NomPrenom string `json:"nomprenom"`
	DateNaissance string `json:"datenaissance"`
	Cin  string `json:"cin"`
	MedCons  string `json:"medcons"`
}

type Receptionniste struct {
	NomPrenom string `json:"nomprenom"`
	Cin  string `json:"cin"`
	Hopital  string `json:"hopital"`
}

type AgentLab struct {
	NomPrenom string `json:"nomprenom"`
	Cin  string `json:"cin"`
	Laboratoire  string `json:"laboratoire"`
}

type Dossier struct {
	CinPatient string `json:"cinpatient"`
	Dossier string `json:"dossier"`
}

type Permission struct {
	CinMedecin string `json:"cinmedecin"`
	CinPatient string `json:"cinpatient"`
	NomPrenomMedecin string `json:"nomprenommedecin"`
	NomPrenomPatient string `json:"nomprenompatient"`
	Specialite  string `json:"specialite"`
	Debautorisation string `json:"debautorisation"`
	Finautorisation string `json:"finautorisation"`
	Statut int `json:"statut"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "tuna-chaincode"
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} 

	if function == "get_data" {
		return s.get_data(APIstub, args)
	} 

	if function == "req_access" {
		return s.req_access(APIstub, args)
	} 

	if function == "fetch_access" {
		return s.fetch_access(APIstub, args)
	} 

	if function == "fetch_perm" {
		return s.fetch_perm(APIstub, args)
	} 

	if function == "pat_fetch_access" {
		return s.pat_fetch_access(APIstub, args)
	} 

	if function == "queryAllPatient" {
		return s.queryAllPatient(APIstub)
	} 

	if function == "edit_perm" {
		return s.edit_perm(APIstub, args)
	} 

	if function == "get_dossier" {
		return s.get_dossier(APIstub, args)
	}


	if function == "add_to_dossier" {
		return s.add_to_dossier(APIstub, args)
	}

	if function == "add_patient" {
		return s.add_patient(APIstub, args)
	}

	if function == "add_medecin" {
		return s.add_medecin(APIstub, args)
	}

	if function == "add_agent" {
		return s.add_agent(APIstub, args)
	}

	return shim.Error(function)
}

/* Get Data */
func (s *SmartContract) get_data(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	dataAsBytes, _ := APIstub.GetState(args[0])
	if dataAsBytes == nil {
		return shim.Error("Data not found.")
	}
	return shim.Success(dataAsBytes)
}

/* InitLegdger */

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	/*patient := []Patient{
		Patient{NomPrenom: "Ayman Belhadj", DateNaissance: "05/05/1997", Cin: "11071155"},
		Patient{NomPrenom: "Nour Mziou", DateNaissance: "01/01/1998", Cin: "87654321"},
	}

	medecin := []Medecin{
		Medecin{NomPrenom: "MedAmin Hentati", Specialite: "Radiologie", Cin: "12345678"},
		Medecin{NomPrenom: "Abc Def", Specialite: "Orthopédie", Cin: "88888888"},
	}


	dossier := []Dossier{
		Dossier{CinPatient:"11071155",Dossier:"1620405339425:M:12345678: test dossier;1620405339425:M:12345678: test dossier2"},
		Dossier{CinPatient:"87654321",Dossier:""},
	}

	agentLab := []AgentLab{
		AgentLab{NomPrenom : "Agent 1", Cin:"12341234", Laboratoire:"Labo 1"},
	}*/

	recepts := []Receptionniste{
		Receptionniste{NomPrenom : "Recep 1", Cin:"11111111", Hopital:"Hopital 1"},
	}

	m := 0
	for m < len(recepts) {
		receptsLabAsBytes, _ := json.Marshal(recepts[m])
		APIstub.PutState("R-"+recepts[m].Cin, receptsLabAsBytes)
		fmt.Println("Added", recepts[m])
		m = m + 1
	}

	/*i := 0
	for i < len(patient) {
		patientAsBytes, _ := json.Marshal(patient[i])
		APIstub.PutState("P-"+patient[i].Cin, patientAsBytes)
		fmt.Println("Added", patient[i])
		i = i + 1
	}

	j := 0
	for j < len(medecin) {
		medecinAsBytes, _ := json.Marshal(medecin[j])
		APIstub.PutState("M-"+medecin[j].Cin, medecinAsBytes)
		fmt.Println("Added", medecin[j])
		j = j + 1
	}

	k := 0
	for k < len(dossier) {
		dossierAsBytes, _ := json.Marshal(dossier[k])
		APIstub.PutState("D-"+dossier[k].CinPatient, dossierAsBytes)
		fmt.Println("Added", dossier[k])
		k = k + 1
	}

	l := 0
	for l < len(agentLab) {
		agentLabAsBytes, _ := json.Marshal(agentLab[l])
		APIstub.PutState("A-"+agentLab[l].Cin, agentLabAsBytes)
		fmt.Println("Added", agentLab[l])
		l = l + 1
	}*/

	

	return shim.Success(nil)
}


/* Get Dossier */
func (s *SmartContract) get_dossier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	dosAsBytes, _ := APIstub.GetState("D-"+args[1])
	if dosAsBytes == nil {
		return shim.Error("Could not locate dossier")
	}

	if args[3] == "PRIVATE"{
		return shim.Success(dosAsBytes)
	}

	indexName := "med-pat"
	Key, err := APIstub.CreateCompositeKey(indexName, []string{args[0], args[1], args[2]})
	if err != nil {
		return shim.Error(err.Error())
	}
	permAsBytes, _ := APIstub.GetState(Key)
	if permAsBytes == nil {
		return shim.Error("Permission not found.")
	}
	permission := Permission{}
	json.Unmarshal(permAsBytes, &permission)

	if permission.Statut < 2{
		return shim.Error("You don't have permission.")
	}

	now := time.Now()
	msec := now.UnixNano() / int64(time.Millisecond)
	finAut, err2 := strconv.ParseInt(permission.Finautorisation, 10, 64)
	if err2 != nil {
		return shim.Error(err2.Error())
	}
	if msec > finAut{
		return shim.Error("You don't have permission.")
	}

	return shim.Success(dosAsBytes)
}


/* Request Access */

func (s *SmartContract) req_access(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}


	medAsBytes, _ := APIstub.GetState("M-"+args[0])
	if medAsBytes == nil {
		return shim.Error("Could not locate medecin")
	}
	medecin := Medecin{}


	json.Unmarshal(medAsBytes, &medecin)

	patAsBytes, _ := APIstub.GetState("P-"+args[1])
	if patAsBytes == nil {
		return shim.Error("Could not locate patient")
	}
	patient := Patient{}


	json.Unmarshal(patAsBytes, &patient)
	
	if !strings.Contains(patient.MedCons,args[0]){
		if len(patient.MedCons)==0 {
			patient.MedCons = fmt.Sprintf("%s", args[0])
		} else {
		patient.MedCons = fmt.Sprintf("%s;%s",patient.MedCons, args[0])
		}
		patAsBytes, _ = json.Marshal(patient)
		APIstub.PutState("P-"+args[1], patAsBytes)
	}
	

	// Normally check that the specified argument is a valid holder of tuna
	// we are skipping this check for this example
	

	var perm = Permission{Specialite :medecin.Specialite, CinMedecin: args[0], CinPatient: args[1], NomPrenomMedecin : medecin.NomPrenom,NomPrenomPatient : patient.NomPrenom, Debautorisation: args[2], Finautorisation: args[3], Statut: 1 }
	indexName := "med-pat"
	Key, err := APIstub.CreateCompositeKey(indexName, []string{args[0], args[1], args[2]})
	if err != nil {
		return shim.Error(err.Error())
	}
	permAsBytes, _ := json.Marshal(perm)
	APIstub.PutState(Key, permAsBytes)

	return shim.Success(nil)
}

/* Fetch Medecin access */

func (s *SmartContract) fetch_access(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	permResultsIterator, err := APIstub.GetStateByPartialCompositeKey("med-pat", []string{args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer permResultsIterator.Close()

	// Iterate through result set and for each marble found, transfer to newOwner
	var i int

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for i = 0; permResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble name from the composite key
		responseRange, err := permResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// get the color and name from color~name composite key
		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		
		idMedecin := compositeKeyParts[0]
		idPatient := compositeKeyParts[1]
		dateCons := compositeKeyParts[2]

		fmt.Printf("- found a marble from index:%s color:%s name:%s\n", objectType, idMedecin, idPatient)


		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		/*str := string(responseRange.Value)
		res := Permission{}
		json.Unmarshal([]byte(str), &res)
		now := time.Now()
		msec := now.UnixNano() / int64(time.Millisecond)
		s := strconv.FormatInt(msec, 10)*/
		buffer.WriteString(idMedecin + "-" + idPatient + "-" + dateCons)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(responseRange.Value))
		
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/* Patient Fetch Access */
func (s *SmartContract) pat_fetch_access(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	

	permResultsIterator, err := APIstub.GetStateByPartialCompositeKey("med-pat", []string{args[0],args[1]})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer permResultsIterator.Close()

	// Iterate through result set and for each marble found, transfer to newOwner
	var i int

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for i = 0; permResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble name from the composite key
		responseRange, err := permResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// get the color and name from color~name composite key
		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		
		idMedecin := compositeKeyParts[0]
		idPatient := compositeKeyParts[1]
		dateCons := compositeKeyParts[2]

		fmt.Printf("- found a marble from index:%s color:%s name:%s\n", objectType, idMedecin, idPatient)

		str := string(responseRange.Value)
		res := Permission{}
		json.Unmarshal([]byte(str), &res)
		now := time.Now()
		msec := now.UnixNano() / int64(time.Millisecond)
		msec30, err := strconv.ParseInt(res.Debautorisation, 10, 64)
		
		if err!=nil{
			shim.Error("ERROR")
		}
		
		if msec < msec30 || args[2] == "ALL" {
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(idMedecin + "-" + idPatient + "-" + dateCons)
			buffer.WriteString("\"")
			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(responseRange.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}
	}

	buffer.WriteString("]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/* Fetch Permission */

func (s *SmartContract) fetch_perm(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	indexName := "med-pat"
	Key, err := APIstub.CreateCompositeKey(indexName, []string{args[0], args[1]})
	if err != nil {
		return shim.Error(err.Error())
	}
	permAsBytes, _ := APIstub.GetState(Key)
	if permAsBytes == nil {
		return shim.Error("Permission not found.")
	}
	return shim.Success(permAsBytes)
}

/* Query All Patient */

func (s *SmartContract) queryAllPatient(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "P-"
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/* Query All Medecin */

func (s *SmartContract) queryAllMedecin(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "M-"
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


/* Edit permission */
func (s *SmartContract) edit_perm(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	indexName := "med-pat"
	Key, err := APIstub.CreateCompositeKey(indexName, []string{args[0], args[1], args[2]})
	if err != nil {
		return shim.Error(err.Error())
	}
	permAsBytes, _ := APIstub.GetState(Key)
	if permAsBytes == nil {
		return shim.Error("Permission not found.")
	}

	permission := Permission{}

	json.Unmarshal(permAsBytes, &permission)
	permission.Statut, err= strconv.Atoi(args[3])
	if args[3] == "2"{
		now := time.Now()
		next30 := now.Add(time.Minute * 30)
		msec := next30.UnixNano() / int64(time.Millisecond)
		s := strconv.FormatInt(msec, 10)
		permission.Finautorisation = s;
	}

	permAsBytes, _ = json.Marshal(permission)
	APIstub.PutState(Key, permAsBytes)
	
	
	return shim.Success(nil)
}

/* Add Agent Lab */
func (s *SmartContract) add_agent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	agent := AgentLab{NomPrenom : args[0], Cin:args[2], Laboratoire:args[1]}
	agentAsBytes, _ := json.Marshal(agent)
	APIstub.PutState("A-"+agent.Cin, agentAsBytes)
	
	return shim.Success(nil)
}

/* Add medecin */
func (s *SmartContract) add_medecin(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	medecin := Medecin{NomPrenom: args[0], Specialite: args[1], Cin: args[2]}
	medecinAsBytes, _ := json.Marshal(medecin)
	APIstub.PutState("M-"+medecin.Cin, medecinAsBytes)
	
	return shim.Success(nil)
}

/* Add patient */
func (s *SmartContract) add_patient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	patient := Patient{NomPrenom: args[0], DateNaissance: args[1], Cin: args[2]}
	patientAsBytes, _ := json.Marshal(patient)
	APIstub.PutState("P-"+patient.Cin, patientAsBytes)
	
	dossier := Dossier{CinPatient: args[2], Dossier:""}
	dossierAsBytes, _ := json.Marshal(dossier)
	APIstub.PutState("D-"+patient.Cin, dossierAsBytes)

	
	return shim.Success(nil)
}

/* Add to dossier */
func (s *SmartContract) add_to_dossier(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	
	dosAsBytes, _ := APIstub.GetState("D-"+args[0])
	if dosAsBytes == nil {
		return shim.Error("Dossier not found.")
	}

	dossier := Dossier{}

	json.Unmarshal(dosAsBytes, &dossier)
	
	if len(dossier.Dossier)==0 {
		dossier.Dossier = fmt.Sprintf("%s", args[1])
	} else {
		dossier.Dossier = fmt.Sprintf("%s;%s",dossier.Dossier, args[1])
	}

	dosAsBytes, _ = json.Marshal(dossier)
	APIstub.PutState("D-"+args[0], dosAsBytes)
	
	
	return shim.Success(nil)
}

/////////////////////////////////////////////////////////

/*
 * The queryTuna method *
Used to view the records of one particular tuna
It takes one argument -- the key for the tuna in question
 */
func (s *SmartContract) queryTuna(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	tunaAsBytes, _ := APIstub.GetState(args[0])
	if tunaAsBytes == nil {
		return shim.Error("Could not locate tuna")
	}
	return shim.Success(tunaAsBytes)
}

/*
 * The initLedger method *
 */


/*
 * The recordTuna method *
Fisherman like Sarah would use to record each of her tuna catches. 
This method takes in five arguments (attributes to be saved in the ledger). 
 */
func (s *SmartContract) recordTuna(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var tuna = Tuna{ Vessel: args[1], Location: args[2], Timestamp: args[3], Holder: args[4] }

	tunaAsBytes, _ := json.Marshal(tuna)
	err := APIstub.PutState(args[0], tunaAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record tuna catch: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllTuna method *
allows for assessing all the records added to the ledger(all tuna catches)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllTuna(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTuna:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The changeTunaHolder method *
The data in the world state can be updated with who has possession. 
This function takes in 2 arguments, tuna id and new holder name. 
 */
func (s *SmartContract) changeTunaHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	tunaAsBytes, _ := APIstub.GetState(args[0])
	if tunaAsBytes == nil {
		return shim.Error("Could not locate tuna")
	}
	tuna := Tuna{}

	json.Unmarshal(tunaAsBytes, &tuna)
	// Normally check that the specified argument is a valid holder of tuna
	// we are skipping this check for this example
	tuna.Holder = args[1]

	tunaAsBytes, _ = json.Marshal(tuna)
	err := APIstub.PutState(args[0], tunaAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change tuna holder: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}