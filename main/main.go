package main

import (
	"time"

	"github.com/doptime/DualModelIterativeReasoning/batchop"
	"github.com/doptime/DualModelIterativeReasoning/message"
	"github.com/doptime/DualModelIterativeReasoning/models"
	"github.com/doptime/DualModelIterativeReasoning/query"
)

// Perform reasoning
var MCTSTrajectory = &query.Query{
	Created: time.Now().Unix(),
	Model:   models.ModelQwen32B.Name,
	MsgUser: message.UserMsg(`I have a 9 yrs old daughter, I want's help here with her using a funny | interesting | breath taking | deep-diving | emotion arousing story. 
Remember, The Most Important thing is building the experience. If can not, Others Fade Away, because she's somehow formidable with her work.
Learn meters, decimeters, and centimeters through a perceptual way. Encode it in a 6,000-word suspense detective novel.
This requires conceiving the outline of the novel, with a lot of twists and turns. There should be 30+ applications of meters, decimeters, and centimeters in total.
The topics include detective, suspense, horror, secret love, elementary school students, experiments, and docks.
Clues: 1.8m fence, 20cm lowered, 15cm scratched, 70cm footprint, 10cm scale, 12cm scale, 500ml beaker, 480ml actual, 2m deep pool, 3 days soaked, 60cm specimen, 1.5m high cabinet, 0.01mm particle, 90cm fracture, 100 boxes, 1m high box, 50kg weight, 1/100 reduction, 1cm cube, 5cm miniature man, 0.5mm DNA, 150cm tall, 148cm tall, 10cm device, 5kg weight loss, 5mm growth, 7mm pupil, 90cm cube, 2cm insect, 30cm sinking, 2m high robot, 2cm lengthened, 0.1% volume, 0.01mm error, 3mm clone, 1km2 city, 40cm reduction
Analyze & very impressive visual depict of each clue needed.Third person perspective。
`)}

func main() {
	// err := MCTSTrajectory.Solute()
	// if err == nil {
	// 	batchop.CopyToClipboard(MCTSTrajectory)
	// }
	// var KeyTreeNode = db.HashKey[string, *Query]()
	// mp, _ := query.KeyTreeNode.HGetAll()

	// for _, v := range mp {
	// 	query.NodesMap.Set(v.Id, v)
	// }
	// query.NodesMap.Set("root", MCTSTrajectory)
	// if node, ok := query.NodesMap.Get("root"); ok {
	// 	MCTSTrajectory = node
	// }
	// if query.NodesMap.Count() == 0 {
	// 	query.NodesMap.Set("root", MCTSTrajectory)
	// }
	problemReformulated, err := batchop.ProblemReformulation(MCTSTrajectory)
	if err != nil {
		return
	}
	msgBest := batchop.SelectBestNode(problemReformulated...)

	batchop.ProblemSolving(msgBest)

}
