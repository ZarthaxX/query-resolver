# TODOs
- Allow flexible logic formulas, transforming them later into normal form.
  - We won't be able to preserve the NOT expression, as it would make Datasource job harder to detect cases where a NOT expression wraps another one, and as a consequence avoid bugs because we forgot about this context when visiting the expression.
- Determine the type of the whole expression, maybe requiring the interface of the different expressions to have a Type() method.