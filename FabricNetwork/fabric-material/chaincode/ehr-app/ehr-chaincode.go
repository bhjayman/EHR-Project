
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	DateNaissance string `json:"datenaissance"`
	Specialite  string `json:"specialite"`
	Grade  string `json:"grade"`
	Cin  string `json:"cin"`
}

type Patient struct {
	NomPrenom string `json:"nomprenom"`
	DateNaissance string `json:"datenaissance"`
	Cin  string `json:"cin"`
}

type Receptionniste struct {
	NomPrenom string `json:"nomprenom"`
	DateNaissance string `json:"datenaissance"`
	Cin  string `json:"cin"`
	Hopital  string `json:"hopital"`
}

type AgentLab struct {
	NomPrenom string `json:"nomprenom"`
	DateNaissance string `json:"datenaissance"`
	Cin  string `json:"cin"`
	Laboratoire  string `json:"laboratoire"`
}

type Dossier struct {
	NumDossier string `json:"numdossier"`
	CinPatient string `json:"cinpatient"`
}

type Permission struct {
	CinMedecin string `json:"cinmedecin"`
	CinPatient string `json:"cinpatient"`
	Debautorisation string `json:"debautorisation"`
	Finautorisation string `json:"finautorisation"`
	Autorise bool `json:"autorise"`
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

	if function == "checkUser" {
		return s.checkUser(APIstub, args)
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

	if function == "queryAllPatient" {
		return s.queryAllPatient(APIstub)
	} 
	/*
	if function == "queryPatient" {
		return s.queryPatient(APIstub, args)
	} 
	if function == "recordPatient" {
		return s.recordPatient(APIstub, args)
	}
	if function == "queryAllPatient" {
		return s.queryAllPatient(APIstub)
	} 
	if function == "updateDossierMedical" {
		return s.updateDossierMedical(APIstub, args)
	}*/
	// query medecin query agent lab record medecin record agent lab .....

	return shim.Error(function)
}

/* Check User */
func (s *SmartContract) checkUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	userAsBytes, _ := APIstub.GetState(args[0])
	if userAsBytes == nil {
		return shim.Error("User not found.")
	}
	return shim.Success(userAsBytes)
}

/* InitLegdger */

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	patient := []Patient{
		Patient{NomPrenom: "Ayman Belhadj", DateNaissance: "05/05/1997", Cin: "11071155"},
		Patient{NomPrenom: "Nour Mziou", DateNaissance: "01/01/1998", Cin: "87654321"},
	}

	medecin := []Medecin{
		Medecin{NomPrenom: "MedAmin Hentati", DateNaissance: "30apr", Specialite: "Radiologie", Grade: "1", Cin: "12345678"},
	}

	i := 0
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

	return shim.Success(nil)
}


/* Request Access */

func (s *SmartContract) req_access(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var perm = Permission{ CinMedecin: args[0], CinPatient: args[1], Debautorisation: args[2], Finautorisation: args[3], Autorise: false }
	indexName := "med-pat"
	Key, err := APIstub.CreateCompositeKey(indexName, []string{args[0], args[1]})
	if err != nil {
		return shim.Error(err.Error())
	}
	permAsBytes, _ := json.Marshal(perm)
	APIstub.PutState(Key, permAsBytes)
	//APIstub.PutState(args[0]+"-"+args[1], permAsBytes)

	return shim.Success(nil)
}

/* Fetch access */

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

		fmt.Printf("- found a marble from index:%s color:%s name:%s\n", objectType, idMedecin, idPatient)


		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		str := string(responseRange.Value)
		res := Permission{}
		json.Unmarshal([]byte(str), &res)
		now := time.Now()
		msec := now.UnixNano() / int64(time.Millisecond)
		s := strconv.FormatInt(msec, 10)
		buffer.WriteString(idMedecin + "-" + idPatient + "-"+ res.Debautorisation+ "-"+ s)
		buffer.WriteString("\"")

		//buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		//buffer.WriteString(string(queryResponse.Value))
		
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
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