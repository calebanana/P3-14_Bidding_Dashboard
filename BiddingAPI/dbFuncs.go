package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Bid struct {
	StudentID string `bson: "studentID"`
	BidAmt    int32  `bson: "bidAmt"`
	BidStatus string `bson: "bidStatus"`
}

type Class struct {
	ClassCode string `bson: "classCode"`
	ClassBids []Bid  `bson: "classBids"`
}

type Module struct {
	ModuleCode    string        `bson: "moduleCode"`
	ModuleName    string        `bson: "moduleName"`
	ModuleClasses []Class       `bson: "moduleClasses"`
}

type Semester struct {
	SemesterStartDate string   
	SemesterModules   []Module
}

func GetAllBids(inputSemStartDate string, inputClassCode string, inputStudentId string, filtered string) Semester{

	fmt.Println(inputSemStartDate, inputStudentId, inputClassCode)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	biddingDB := client.Database("BiddingDB")
	semCollection := biddingDB.Collection(inputSemStartDate)

	cursor, err := semCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var modules []bson.M
	if err = cursor.All(ctx, &modules); err != nil {
		log.Fatal(err)
	}

	var semesterModules []Module
	for _, i := range modules {
		moduleCode := i["moduleCode"]
		moduleName := i["moduleName"]
		classes := i["moduleClasses"].(primitive.A)

		var moduleClasses []Class
		for _, j := range classes {
			_class := j.(primitive.M)
			classCode := _class["classCode"]
			class_bids := _class["classBids"].(primitive.A)

			var classBids []Bid
			for _, k := range class_bids {
				studentBid := k.(primitive.M)
				studentId := studentBid["studentID"]
				studentBidAmt := studentBid["bidAmt"]
				studentBidStatus := studentBid["bidStatus"]

				var bid = Bid{
					StudentID: studentId.(string),
					BidAmt: studentBidAmt.(int32),
					BidStatus: studentBidStatus.(string),
				}
				if (inputStudentId != ""){
					if bid.StudentID == inputStudentId{
						classBids = append(classBids, bid)
					}
				} else {
					//fmt.Println("INPUT STUDENT ID EMPTY")
					classBids = append(classBids, bid)
				}
			}

			var class = Class{
				ClassCode: classCode.(string),
				ClassBids: classBids,
			}

			if filtered == "true" {
				if (len(class.ClassBids) == 0){
					continue
				}
			}

			if (inputClassCode != "") {
				if class.ClassCode == inputClassCode{
					moduleClasses = append(moduleClasses, class)
				}
			} else {
				moduleClasses = append(moduleClasses, class)
			}
		}

		var module = Module{
			ModuleCode: moduleCode.(string),
			ModuleName: moduleName.(string),
			ModuleClasses: moduleClasses,
		}

		if filtered == "true" {
			if (len(module.ModuleClasses) == 0){
				continue
			}
		}

		if (inputClassCode != ""){
			var inputModuleCode string = inputClassCode[:len(inputClassCode) - 3]
			if module.ModuleCode == inputModuleCode{
				semesterModules = append(semesterModules, module)
			}
		} else {
			semesterModules = append(semesterModules, module)
		}
	}

	var semester = Semester{
		SemesterStartDate: inputSemStartDate,
		SemesterModules: semesterModules,
	}

	return semester
}

func AddNewSemester(inputNewSemester Semester) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	biddingDB := client.Database("BiddingDB")
	newCollection := biddingDB.Collection(inputNewSemester.SemesterStartDate)

	for moduleIndex, module := range inputNewSemester.SemesterModules{
		newModule, _ := newCollection.InsertOne(ctx, bson.D{
			{Key: "moduleCode", Value: module.ModuleCode},
			{Key: "moduleName", Value: module.ModuleName},
			{Key: "moduleClasses", Value: bson.A{
			},
			},
		})

		fmt.Println(newModule.InsertedID)

		for _, class := range inputNewSemester.SemesterModules[moduleIndex].ModuleClasses{
			result, _ := newCollection.UpdateOne(
				ctx,
				bson.M{"moduleCode": module.ModuleCode},
				bson.D{
					{"$push", bson.M{"moduleClasses": bson.M{
						"classCode": class.ClassCode,
						"classBids": bson.A{
						},
					},
					},
					},
				},
			)
		
			fmt.Println("Classes added: ", result.ModifiedCount)
		}
	}
}

func AddNewBid(inputSemStartDate string, inputClassCode string, inputStudentId string, inputBidAmt int32) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	biddingDB := client.Database("BiddingDB")
	semCollection := biddingDB.Collection(inputSemStartDate)

	result, err := semCollection.UpdateOne(
		ctx,
		bson.M{"moduleClasses.classCode": inputClassCode},
		bson.D{
			{"$push", bson.M{"moduleClasses.$.classBids": bson.D{
				{Key: "studentID", Value: inputStudentId},
				{Key: "bidAmt", Value: inputBidAmt},
				{Key: "bidStatus", Value: "Pending"},
			},
			}},
		},
	)

	fmt.Println("Documents added: ", result.ModifiedCount)
}

func EditBid(inputSemStartDate string, inputClassCode string, inputStudentId string, inputBidAmt int32) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	biddingDB := client.Database("BiddingDB")
	semCollection := biddingDB.Collection(inputSemStartDate)

	inputModuleCode := inputClassCode[:len(inputClassCode) - 3]
	filter := bson.D{{"moduleCode", inputModuleCode}}
	update := bson.D{
		{"$set", bson.D{
			{"moduleClasses.$[class].classBids.$[bid].bidAmt", inputBidAmt},
		}},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
        Filters: []interface{}{
			bson.M{"class.classCode": inputClassCode},
			bson.M{"bid.studentID": inputStudentId},
		},
    })

	result, err := semCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Documents modified: ", result.ModifiedCount)
}

func DeleteBid(inputSemStartDate string, inputClassCode string, inputStudentId string) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	biddingDB := client.Database("BiddingDB")
	semCollection := biddingDB.Collection(inputSemStartDate)
	
	inputModuleCode := inputClassCode[:len(inputClassCode) - 3]
	filter := bson.D{{"moduleCode", inputModuleCode}}
	update := bson.D{
		{"$pull", bson.D{
			{"moduleClasses.$[class].classBids", bson.D{
				{"studentID", inputStudentId},
			}},
		}},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
        Filters: []interface{}{
			bson.M{"class.classCode": inputClassCode},
		},
    })

	result, err := semCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Documents deleted: ", result.ModifiedCount)
}