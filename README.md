# **Word of the Day (WOTD) on the Command Line**

## Installation

1. Install Go 1.17+
2. $ ```git clone https://github.com/mcdxwell/wotd.git```
3. $ ```cd wotd```
4. $ ```go build```


## Commands

1. $ ```wotd get```
> Output: $ Cryptography
2. $ ```wotd random```
> Ouput: $ Itinerant

## Known Issues

- Merriam-Webster updates the word of the day around 9:00 AM UTC, so if you do `wotd get` before the wotd is updated, then the wotd will be repeated inside the json file. 
- Quick fix: manually change the date of the word of the day in the wotd.json file.


## Notes

- Update the path in words.go to save words to the wotd.json in a desired location.
- This allows for all words to be saved in one place when `wotd get` or `wotd random` is used globally.