### Evaluate the value of an arithmetic expression in Reverse Polish (Postfix) Notation

```
The input will be a postfix string and for ease the character will be separated by a space. The function should return a Number. Inputs must
consider positive and negative numbers. Polish notation is a way of expressing arithmetic expressions. Its most basic distinguishing feature
is that operators are placed on the left of their operands.
Your implementation can consider operators ^ , * , / , + , - and basic trigonometric functions.
For more information on how to solve RPN read Reverse Polish Notation on wikipedia.

Sample Input 1:
10 3 +
(Traditionally, 10 + 3)
Output 1 :
13

Sample Input 2:
10 3 2 + -
(Traditionally, 10 - ( 3 + 2) )
Output 2:
5

Sample Input 3:
1 10 3 * 2 ^
(Traditionally, (10 * 3)^2 )
Output 3:
900
```

### How to use

~~~
go run RPNCalculator.go
~~~

`Supported operations`: 

```
+, -, *, /, ^, sin, cos, tan, asin, acos, atan, sqrt, ctg
```

`Example`:

~~~ 
Enter an expression:
10 3 +
Result: 13
Enter an expression:
10 4 /
Result: 2.5
Enter an expression:
1 sin
Result: 0.841
Enter an expression:
4 16 sqrt sqrt +
Result: 6
Enter an expression:
+
not enough elements to calculate
4 5 + +
not enough elements to calculate
4 0 /
can not divide by zero
abc
token is invalid: abc
5
Result: 5
~~~

One can find tests in RPNCalculator_test.go.
~~~
go test
~~~