# **Word of the Day (wotd) on the Command Line**

## Installation

1. Install Go 1.21+
1. $ ```git clone https://github.com/mcdxwell/wotd.git```
1. $ ```cd wotd```
1. $ ```go build .``` or ```go install .```

## Commands

1. $ ```wotd```
> Output: $ Cryptography
1. $ ```wotd random```
> Ouput: $ Itinerant

## Known Issues

- Merriam-Webster updates the word of the day around 9:00 AM UTC, so if you do `wotd` before the wotd is updated, then the wotd will be repeated inside the json file. 
- Quick fix: manually change the date of the word of the day in the wotd.json file.

## Notes

- Update the path in words.go to save words to the wotd.json in a desired location.
- This allows for all words to be saved in one place when `wotd` or `wotd random` is used globally.