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

func GetAllBids(inputSemStartDate string, inputClassCode string, inputStudentId string) Semester{

	fmt.Println( inputSemStartDate, inputStudentId, inputClassCode)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://bidding_db:27017"))
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

			if (len(class.ClassBids) == 0){
				continue
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
		if (inputClassCode != ""){
			var inputModuleCode string = inputClassCode[:len(inputClassCode) - 2]
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

func AddNewSemester(inputSemStartDate string) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://bidding_db:27017"))
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
	newCollection := biddingDB.Collection(inputSemStartDate)

	adb, _ := newCollection.InsertOne(ctx, bson.D{
		{Key: "moduleCode", Value: "ADB"},
		{Key: "moduleName", Value: "Advanced Databases"},
		{Key: "moduleClasses", Value: bson.A{
			bson.M{
				"classCode": "ADB01",
				"classBids": bson.A{
				},
			},
			bson.M{
				"classCode": "ADB02",
				"classBids": bson.A{
				},
			},
		},
		},
	})

	dl, _ := newCollection.InsertOne(ctx, bson.D{
		{Key: "moduleCode", Value: "DL"},
		{Key: "moduleName", Value: "Deep Learning"},
		{Key: "moduleClasses", Value: bson.A{
			bson.M{
				"classCode": "DL01",
				"classBids": bson.A{
				},
			},
			bson.M{
				"classCode": "DL02",
				"classBids": bson.A{
				},
			},
		},
		},
	})

	fmt.Println(adb.InsertedID)
	fmt.Println(dl.InsertedID)
}

func AddNewBid(inputSemStartDate string, inputClassCode string, inputStudentId string, inputBidAmt int32) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://bidding_db:27017"))
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

	_, err = semCollection.UpdateOne(
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
}

// func EditBid(inputSemStartDate string, inputClassCode string, inputStudentId string, newBidAmt int32) {
// 	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Disconnect(ctx)
// 	biddingDB := client.Database("testdb")
// 	semCollection := biddingDB.Collection(inputSemStartDate)

// 	unwind1 := bson.D{{"$unwind", "moduleClasses"}}
// 	cursor, err := semCollection.Aggregate(ctx, mongo.Pipeline{unwind1})
// 	if err != nil {
//     	panic(err)
// 	}
// 	var data []bson.M
// 	if err = cursor.All(ctx, &data); err != nil {
//     	panic(data)
// 	}
// 	fmt.Println(showsWithInfo)

// 	// _, err = semCollection.UpdateOne(
// 	// 	ctx,
// 	// 	bson.M{"moduleClasses.classCode": inputClassCode, "moduleClasses.classBids.studentID": inputStudentId},
// 	// 	bson.D{
// 	// 		{"$set", bson.M{"moduleClasses.$.classBids.studentID": bson.D{
// 	// 			{Key: "studentID", Value: inputStudentId},
// 	// 			{Key: "bidAmt", Value: newBidAmt},
// 	// 			{Key: "bidStatus", Value: "Pending"},
// 	// 		},
// 	// 		}},
// 	// 	},
// 	// )
// }

// func DeleteBid(inputSemStartDate string, inputClassCode string, inputStudentId string, newBidAmt int32) {
// 	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Disconnect(ctx)
// 	biddingDB := client.Database("testdb")
// 	semCollection := biddingDB.Collection(inputSemStartDate)

	
// }