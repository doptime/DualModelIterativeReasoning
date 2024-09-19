This Project aims to realize openAI o1 reasoning

Methodology:
1. Question reasoing & rephrasing
   - In most case, a question can not be well-solved, becasue the model can not find out the real meaning/intension of the question. So the LLM tring to solve question literally, ability can't be fully triggered.
   - Try rephrasing the question, and until the boundary condition, restriction, context, and destination is converged.
2. Planing & Solving in Iterative way
   - Planing & Solving are try to carried as soon as possible to close the improving loop . The real important thing is End2End Iterating, review the planning and solution, and iterate.
3. Mutual Reasoning
   - As rStar paper shows, self-iterating and correcting is not effective, and multi-model is necessary to boost the performance.
   - Use Multi-Model or Dual Modal to cross-review the logic chain. rStar paper shows that self-iterating and correcting is not effective, and multi-model is necessary to boost the performance.
4. Parallel-beam searching
   - Use parallel-beam searching to balance the exploration and exploitation. Single instance may have performance degradation during exploration. Reasoning improvement can not rely on single branch iteration.