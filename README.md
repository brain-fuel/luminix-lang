# Acorn Language

## Rules of Engagement

- Engage sincerely
- Think deeply
- Build something worth building

## Ideas

- propositional and predicate logic
- `if P then Q -> P => Q`; `P => Q -> if P then Q`... make a facility to interchange one of these for another. Perhaps also if we have something like `assert P => Q ;; observe ~P ;; Q: unknown`. We need to have different verbs like `assert`, `observe`, `conclude`, etc. Perhaps we could integrate the Answer Set solver here.
- make an autoformatter that takes in raw files and converts to either math or english form
- make an autoformatter that fixes REPL history
- permit `neither x nor y` verbiage
- `:test` command from REPL to run tests asynchronously and print basic results and way to follow them once they have resolved

## Roadmap

## lang0.1.1

### repl

- carat showing position of actual error rather than arbitrarily at beginning of line
- type of results
- painless navigation
- helps for default functions

## Features

- lx repl (Read, Evaluate, Print Loop): this allows for quick exploration and experimentation in a realm
- boolean logic

# TODO

## Definition of done
For any of the boolean things:
- tests for everything added
- if applicable, english to math
- if applicable, math to english
- seamless repl integration (this means backspace and switching work appropriately)

## lang/parser

### boolean
#### Binary Operators
- and `/\`
- nand `~/\`
- or `\/`
- nor `~\/`
- xnor iff `<=>`
- implies `=>`
- is implied by `<=`
- inhibits `/=>`
- is inhibited by `<=/`
- left `<s`
- right `s>` 
- not left `</`
- not right `/>`

## lang/parser

## lang/evaluator

### boolean
### Primitives
- true #t
- false #f
### Unary Operators
- not / ~
- nullify
- truify
- id
### Binary Operators
- and /\
- nand ~/\
- or \/
- nor ~\/
- xnor iff <=>
- implies =>
- is implied by <=
- inhibits /=>
- is inhibited by <=/
- left `<s`
- right `s>` 
- not left `</`
- not right `/>`

# DONE

## lang/parser

### boolean
- True
- False

#### Unary Operators
- nullify
- truify
- id

## lang/evaluator

### boolean
