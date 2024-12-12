# Inference Engine

## Objective
The goal of this project is to create an inference engine that can be used to infer facts 
from a set of rules and facts. The engine should manage the identity of the concepts and 
uncertainty.
One Possible application is to create a triage system for a hospital, or gamification System. 

## Definitions
A rule is a statement that can be true or false, a fact is a statement that is true.
A Fact can be inferred from a set of rules and facts.
A Fact has a values set by itself depending on the nature of the fact these will vary.
A Rule has a condition and a conclusion, the condition is a set of facts that must be true for the rule to be true.
A Conclusion is a set of facts that can be inferred from the rule.
An Inference can be modifiable or not, other inferences can change a modifiable inference.
A Contradiction is a set of facts that cannot be true at the same time.
The engine should detect contradictions and solve it.

## State of the project
The project is in the early stages of development; the first step is to create the representation

## Representation
the representation of the rules and facts will be done using the expr language.
https://github.com/expr-lang/expr
