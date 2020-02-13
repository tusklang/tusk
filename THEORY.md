# Micro-Calculations
Omm is not like other programming languages. Where other programming languages use binary to calculate, Omm uses, what I like to call, micro-calculations.

# Theory

A micro-calculation is just a calculation using binary for each section of both numbers.

```
192833823 + 273917393
Split -->
192 + 833 + 823, 273 + 917 + 393
Group -->
192 + 273, 833 + 917, 823 + 393
Calculate (Using Binary) -->
465, 1750, 1216
Carry -->
465 + carry(1), 750 + carry(1), 216
Combine -->
466751216
```

This also works for decimals, there is just an extra step: decimalizing

```
8.221 + 2.91
Format -->
8.221 + 2.910
De-Decimalize -->
8221 + 2910
Add (Using Previous) -->
11131
Re-Decimalize (Given the decimal place of the formatted numbers) -->
11.131
```

Each operator has its own rule for calculation.

Micro-calculations are similar to how humans do math

```
  11
  919
+ 281
-------
 1200
```

But it is a bit more complex due to the formatting issues, decimalizations, splitting, and grouping.

# Discussion

Since there are so many more steps in micro-calculations than in traditional programming operations, Omm is relatively slow with small calculations, e.g. 1 + 1, but when it comes to large calculations, Omm is incredibly precise and efficient.

Before you consider using Python or Java to do these large number calculations, consider this:
Python and Java are based on a binary number system. This restricts decimals and, to an extend, integers. Omm's numbers have arbitrary precision, even more than those of Python.
