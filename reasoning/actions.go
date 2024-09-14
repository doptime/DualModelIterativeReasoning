package reasoning

type Action uint32

var A1ProposeOneStepThought Action = 1
var A2ProposeRemainingThoughtSteps Action = 2
var A3ProposeNextSubQuestionAlongWithItsAnswer Action = 4
var A4AnswerSubQuestionAgain Action = 8
var A5RephraseQuestionOrSubQuestion Action = 16
