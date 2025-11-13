# **Word of the Day (wotd) on the Command Line**

## Installation

Install Go 1.21+

`$ git clone https://github.com/mcdxwell/wotd.git`

`$ cd wotd`

`$ go build .` or `$ go install .`


## Known Issues

- Merriam-Webster updates the word of the day around 9:00 AM UTC, so if you do `wotd` before the wotd is updated, then the wotd will be repeated inside the json file. 
- Quick fix: manually change the date of the word of the day in the wotd.json file.
