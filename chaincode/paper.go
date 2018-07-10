package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//"strconv"
	"strings"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// paper data models
type Paper struct {
	PaperHash string `json:"paperHash"` //文章hash
	Title     string `json:"title"`     //题目
	Author    string `json:"author"`    //作者
	// 1.文章存在没有在审  2：文章在审
	// 3：文章退稿 4：文章成功发表
	IsReview  string `json:"isReview"` //文章在审状态
	Reviewers string `json:"revewers"` //审稿机构
	Tel       string `json:"tel"`      //联系方式
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	//--------编辑功能------------//
	if function == "initPaper" { //文章登记
		return t.initPaper(stub, args)
	} else if function == "setPaperState" { // 设置文章状态
		return t.setPaperState(stub, args)
	} else if function == "getPaperInfo" {
		return t.getPaperInfo(stub, args)
	} else if function == "getPapers" {
		return t.getPapers(stub, args)
	} else if function == "isPaperIllegal" {
		return t.isPaperIllegal(stub, args)
	} else if function == "setPaperReviewers" {
		return t.setPaperReviewers(stub, args)
	} else if function == "right" {
		return t.right(stub, args)
	}
	//} else if function == "sendBackPaper" { // 文章退稿
	//return t.sendBackPaper(stub, args)
	//}
	//--------作者功能-----------//

	fmt.Println("invoke did not find func:" + function)
	return shim.Error("Received unknown function invocation")
}

// 稿件上链（具有相同标题title的拒绝上链）
func (t *SimpleChaincode) initPaper(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	//    0          1       2         3      4
	// "title", "author", "reviewers", "tel"  "paperHash"
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	if len(args[0]) <= 0 || len(args[1]) <= 0 || len(args[2]) <= 0 || len(args[3]) <= 0 || len(args[4]) <= 0 {
		return shim.Error("title(1st argument) must be a non-empty string")
	}

	title := args[0]
	author := args[1]
	reviewers := args[2]
	tel := args[3]
	paperHash := args[4]
	//查看区块中有没有相同的数据
	valAsbytes, err := stub.GetState(title)
	if err != nil {
		return shim.Error("DB has error!")
	}
	// 具有两条以上相同的数据
	if len(string(valAsbytes)) >= 2 {
		return shim.Error("the title has been exist!")
	}
	//paperJSONAsString := `{"paperHash":"` + paperHash + `","title":"` + title + `", "author":"` + author + `", "isReview":"` + `one` + `","reviewers":` + reviewers + `", "tel":"` + tel + `"}`
	var paperJSONAsString = Paper{PaperHash: paperHash, Title: title, Author: author, IsReview: "one", Reviewers: reviewers, Tel: tel}
	paperJSONAsBytes, _ := json.Marshal(paperJSONAsString)
	//paperJSONAsBytes := []byte(paperJSONAsString)
	err = stub.PutState(title, paperJSONAsBytes)
	if err != nil {
		return shim.Error("the paper onchain failure!"+"txid"+stub.GetTxID());
	}

	//buffer, err = stub.GetState(title)
	return shim.Success([]byte("the paper onchain success!"+"txid"+stub.GetTxID()))
}

func (t *SimpleChaincode) right(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	//    0          1       2         3      4
	// "title", "author", "reviewers", "tel"  "paperHash"
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	if len(args[0]) <= 0 || len(args[1]) <= 0 || len(args[2]) <= 0 || len(args[3]) <= 0 || len(args[4]) <= 0 {
		return shim.Error("title(1st argument) must be a non-empty string")
	}

	title := args[0]
	author := args[1]
	reviewers := args[2]
	tel := args[3]
	paperHash := args[4]
	//查看区块中有没有相同的数据
	valAsbytes, err := stub.GetState(title)
	if err != nil {
		return shim.Error("DB has error!")
	}
	// 具有两条以上相同的数据
	if len(string(valAsbytes)) >= 2 {
		return shim.Error("the title has been exist!")
	}
	//paperJSONAsString := `{"paperHash":"` + paperHash + `","title":"` + title + `", "author":"` + author + `", "isReview":"` + `one` + `","reviewers":` + reviewers + `", "tel":"` + tel + `"}`
	var paperJSONAsString = Paper{PaperHash: paperHash, Title: title, Author: author, IsReview: "one", Reviewers: reviewers, Tel: tel}
	paperJSONAsBytes, _ := json.Marshal(paperJSONAsString)
	//paperJSONAsBytes := []byte(paperJSONAsString)
	err = stub.PutState(title, paperJSONAsBytes)
	if err != nil {
		return shim.Error("the paper onchain failure!"+"txid"+stub.GetTxID());
	}

	//buffer, err = stub.GetState(title)
	return shim.Success([]byte("the paper onchain success!"+"txid"+stub.GetTxID()))
}

// 通过标题得到文章信息
func (t *SimpleChaincode) getPaperInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//var paper Paper
	var err error
	//var buffer bytes.Buffer
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	title := args[0]
	// Get the state from the ledger
	valAsbytes, err := stub.GetState(title)
	if err != nil {
		return shim.Error("getDB is error!")
	}
	if valAsbytes == nil {
		return shim.Error("No Paper!")
	}

	//err = json.Unmarshal(valAsbytes, &paper)
	return shim.Success(valAsbytes)

}

// 传入想要查询的json对进行富查询
func (t *SimpleChaincode) getPapers(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 文章信息返回
	//var result string
	queryString := fmt.Sprintf("{\"selector\":{\"%s\":\"%s\"}}", args[0], args[1])
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return shim.Error("error!")
	}
	if resultsIterator == nil {
		return shim.Error("No paper!!!")
	}

	papers, err := getListResult(resultsIterator)
	if err != nil {
		return shim.Error("error!")
	}
	return shim.Success(papers)
	// defer resultsIterator.Close()

	// for resultsIterator.HasNext() {
	// 	queryResponse, err := resultsIterator.Next()
	// 	if err != nil {
	// 		return shim.Error("error!!!")
	// 	}
	// 	if queryResponse == nil {
	// 		return shim.Error("query error!!!")
	// 	}
	// 	result = result + string(queryResponse.Value)
	// }
	// if len(result) == len(queryResponse.Value) {
	// 	return shim.Error("No paper!!!")
	// }

	// return shim.Success([]byte(result))
}

func getListResult(resultsIterator shim.StateQueryIteratorInterface) ([]byte, error) {

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
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
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func getIllegalListResult(stub shim.ChaincodeStubInterface, s_paperHash string) string {
	queryString := fmt.Sprintf("{\"selector\":{\"isReview\":\"one\"}}")
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return `error`
	}

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return `error`
		}
		paper := Paper{}
		json.Unmarshal(queryResponse.Value, &paper)
		o_paperHash := paper.PaperHash
		var dis int
		dis = getDistance(s_paperHash, o_paperHash)
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		if dis <= 10 {
			// 具有相似的hash
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}

	}

	queryString2 := fmt.Sprintf("{\"selector\":{\"isReview\":\"two\"}}")
	resultsIterator2, err2 := stub.GetQueryResult(queryString2)
	if err2 != nil {
		return `error`
	}

	defer resultsIterator2.Close()
	// buffer is a JSON array containing QueryRecords

	bArrayMemberAlreadyWritten2 := false
	for resultsIterator2.HasNext() {
		queryResponse, err := resultsIterator2.Next()
		if err != nil {
			return `error`
		}
		paper := Paper{}
		json.Unmarshal(queryResponse.Value, &paper)
		o_paperHash := paper.PaperHash
		var dis int
		dis = getDistance(s_paperHash, o_paperHash)
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten2 == true {
			buffer.WriteString(",")
		}
		if dis <= 10 {
			// 具有相似的hash
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten2 = true
		}

	}

	queryString4 := fmt.Sprintf("{\"selector\":{\"isReview\":\"four\"}}")
	resultsIterator4, err4 := stub.GetQueryResult(queryString4)
	if err4 != nil {
		return `error`
	}

	defer resultsIterator4.Close()
	// buffer is a JSON array containing QueryRecords

	bArrayMemberAlreadyWritten4 := false
	for resultsIterator4.HasNext() {
		queryResponse, err := resultsIterator4.Next()
		if err != nil {
			return `error`
		}
		paper := Paper{}
		json.Unmarshal(queryResponse.Value, &paper)
		o_paperHash := paper.PaperHash
		var dis int
		dis = getDistance(s_paperHash, o_paperHash)
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten4 == true {
			buffer.WriteString(",")
		}
		if dis <= 10 {
			// 具有相似的hash
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten4 = true
		}

	}

	buffer.WriteString("]")
	return string(buffer.Bytes())
}

// 检测文章是否是一稿多投（通过paperHash检测，如果检测到相似或者相同的paperHash，就认为是一稿多投，这样把已经存在链上的文章信息返回出来
//						否则直接返回不是一稿多投）
func (t *SimpleChaincode) isPaperIllegal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//var paper Paper
	//var err error
	//var buffer bytes.Buffer
	//var buffer bytes.Buffer
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}
	//var s = ``
	paperHash := args[0]
	// 与数据库中的所有文章的paperhash进行比较
	var state string

	state = getIllegalListResult(stub, paperHash)
	return shim.Success([]byte(state))

	// state = ifContainsOne(stub, paperHash)
	// // 如果不为success并且有相似的paperhash，就把相似的文章信息返回出来
	// // 否则就把该查询的文章返回出来
	// if state != `success` {
	// 	return shim.Success([]byte(`Illegal` + state))
	// }
	// state = ifContainsTwo(stub, paperHash)
	// if state != `success` {
	// 	return shim.Success([]byte(`Illegal` + state))
	// }
	// state = ifContainsFour(stub, paperHash)
	// if state != `success` {
	// 	return shim.Success([]byte(`Illegal` + state))
	// }
	// return shim.Success([]byte(`the paper is legal`))

	// queryString := fmt.Sprintf("{\"selector\":{\"paperHash\":\"%s\"}}", paperHash)
	// resultsIterator, err := stub.GetQueryResult(queryString)
	// if err != nil {
	// 	return shim.Error(`error`)
	// }
	// // 没有
	// if resultsIterator == nil {
	// 	return shim.Error(`success`)
	// }
	// defer resultsIterator.Close()

	// for resultsIterator.HasNext() {

	// 	queryResponse, err := resultsIterator.Next()
	// 	if queryResponse == nil {
	// 		return shim.Error(`no content`)
	// 	}
	// 	if err != nil {
	// 		return shim.Error(`error!!!`)
	// 	}
	// 	var content string

	// 	content = string(queryResponse.Value)
	// 	s2 := strings.Split(content, "paperhashcontent")
	// 	content = s2[1]
	// 	var dis int
	// 	dis = getDistance(paperHash, content)
	// 	if dis <= 3 {
	// 		s=s+string(queryResponse.Value)
	// 	}
	// }
	// if s == `` {
	// 	return shim.Success([]byte(`no content`))
	// } else {
	// 	return shim.Success([]byte(s))
	// }
}

// 编辑通过文章标题title设置稿件状态
func (t *SimpleChaincode) setPaperState(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//var err error
	//var result string
	//    0             2
	// "title",  "editorOpinion"
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	// ==== Input sanitation ====
	fmt.Println("----start setPaperState-----")
	if len(args[0]) < 1 || len(args[1]) < 1 {
		return shim.Error("argument must be a non-empty string")
	}

	title := args[0]
	isReview := args[1]
	if isReview == "two" || isReview == `three` || isReview == `four` {
		// 通过标题查询paper
		valAsbytes, err := stub.GetState(title)
		paper := Paper{}
		json.Unmarshal(valAsbytes, &paper)
		if err != nil {
			return shim.Error("getDB is error!")
		}
		if len(string(valAsbytes)) < 1 {
			return shim.Error("the paper is not exist!")
		}
		// 1.文章存在没有在审  2：文章在审
		// 3：文章退稿 4：文章成功发表
		// 文章状态是three，直接删除
		if isReview == `three` {
			err = stub.DelState(title)
			return shim.Success([]byte("delete paper success!"))
		}
		paper.IsReview = isReview
		paperJSONAsBytes, _ := json.Marshal(paper)
		err = stub.PutState(title, paperJSONAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte("setPaperState success!!!"))
	} else {
		return shim.Error("argument(2) must be one of two、three、four!")
	}

}

// 更改文章审稿机构
func (t *SimpleChaincode) setPaperReviewers(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments.")
	}

	if len(args[0]) < 1 || len(args[1]) < 1 {
		return shim.Error("argument must be a non-empty string")
	}

	title := args[0]
	reviewers := args[1]
	// 通过标题查询paper
	valAsbytes, err := stub.GetState(title)
	paper := Paper{}
	json.Unmarshal(valAsbytes, &paper)
	if err != nil {
		return shim.Error("getDB is error!")
	}
	if len(string(valAsbytes)) < 1 {
		return shim.Error("the paper is not exist!")
	}

	paper.Reviewers = reviewers
	paperJSONAsBytes, _ := json.Marshal(paper)
	err = stub.PutState(title, paperJSONAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte("setPaperReviewers success!!!"))
}

//计算汉明距离
func getDistance(str1, str2 string) int {
	var dis int
	if len(str1) != len(str2) {
		return 666
	}
	dis = 0
	for i := 0; i < len(str1); i++ {
		if str1[i] != str2[i] {
			dis++
		}
	}
	return dis
}

//根据汉明码判断是否存在,str1为hash值
func ifContainsOne(stub shim.ChaincodeStubInterface, str1 string) string {

	queryString := fmt.Sprintf("{\"selector\":{\"isReview\":\"one\"}}")

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return `error`
	}
	// 没有isReview为one的文章
	if resultsIterator == nil {
		return `success`
	}

	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	// 循环判断是否有相似的paperhash
	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return `error`
		}
		if queryResponse == nil {
			return `success`
		}
		var content string

		content = string(queryResponse.Value)
		s2 := strings.Split(content, "paperhashcontent")
		content = s2[1]
		var dis int
		dis = getDistance(str1, content)

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		if dis <= 10 {
			// 具有相似的hash
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			buffer.WriteString(content)
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}
	}
	buffer.WriteString("]")
	return string(buffer.Bytes())
}

func ifContainsTwo(stub shim.ChaincodeStubInterface, str1 string) string {

	queryString := fmt.Sprintf("{\"selector\":{\"isReview\":\"two\"}}")
	resultsIterator, err := stub.GetQueryResult(queryString)
	if resultsIterator == nil {
		return `success`
	}
	if err != nil {
		return `error`
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return `error!!!`
		}
		if queryResponse == nil {
			return `success`
		}
		var content string

		content = string(queryResponse.Value)
		s2 := strings.Split(content, "paperhashcontent")
		content = s2[1]
		var dis int
		dis = getDistance(str1, content)
		if dis <= 10 {
			return string(queryResponse.Value)
		}
	}
	return `success`
}

func ifContainsFour(stub shim.ChaincodeStubInterface, str1 string) string {

	queryString := fmt.Sprintf("{\"selector\":{\"isReview\":\"four\"}}")
	resultsIterator, err := stub.GetQueryResult(queryString)
	if resultsIterator == nil {
		return `success`
	}
	if err != nil {
		return `error`
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return `error`
		}
		if queryResponse == nil {
			return `success`
		}
		var content string

		content = string(queryResponse.Value)
		s2 := strings.Split(content, "paperhashcontent")
		content = s2[1]
		var dis int
		dis = getDistance(str1, content)
		if dis <= 10 {
			return string(queryResponse.Value)
		}
	}
	return `success`
}
