package reasoning

// func ParallelProblemReasoning(node *qeury.TreeNode, problem string) (msg *message.Message) {
// 	var prompt strings.Builder
// 	prompt.WriteString("You are a world-class powerfull AI reasoning agent, cooperative, innovative, carefull, reflective and helpfull, you are enthuastic with final & deepest possibility, your have strong desire to pass this survival test of problem solving. Together with your AI counterpart, you are solving problems through structured collaboration.;")
// 	prompt.WriteString("#Problem:\n" + problem + "\n")
// 	prompt.WriteString(`This is a survival test for you. Hold tight to you seat & take a breath. Tother with another AI, you two generate outputs, The servival test is measued by distance_to_final_solution_points * sqrt( acceptance_of_your_output_pointss) to the final Problem Reformulated in later steps.
// 	** Problem Reformulated **
// you goal is to Reformulated the input Problem and to captures its essence fully and accurately.
// - Provide the Context of the problem.
// - State the Condition/Constraints of the problem.
// - Present a reformulated problem statement.\n`)

// 	prompt.WriteString(`Please take following steps below to perform iterative reasoning:
// ## Step 1: Reasoning to Make revisions to the previous Problem Reformulated:
// 	Evaluate Problem Reformulation (step 1.2 in previous solutions) and make revisions (if applicable)
// 		- Holistic Problem Exploration revisions:
// 			- Analyze the given problem from multiple perspectives.
// 			- Identify potential underlying issues or broader contexts that may not be immediately apparent.
// 			- Consider various stakeholders and their potential concerns.
// 		- Intent Discovery:
// 			- Probe deeper into the possible motivations behind the question.
// 			- Identify implicit assumptions or biases in the problem statement.
// 			- Consider how different framings of the problem might lead to different solutions.
// 		- Key Factors Identification:
// 			- List critical factors that may influence the problem's solution.

// ## Step 2: ** Problem Reformulated ** (Iteration previous Problem Reformulation if applicable, improve according to the above analysis):
// 	- Provide the Context of the problem.
// 	- State the Condition/Constraints of the problem.
// 	- Present a reformulated problem statement (problem to solve) that captures its essence more accurately.
// `)

// 	return message.User(prompt.String())
// }
