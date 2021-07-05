const express = require('express');
const cors = require("cors");
const fs = require('fs');
const readline = require('readline');

const app = express();
app.use(cors());
app.use(express.json());

const games = [];

const loadJson = async () => {
    try {
        const stream = fs.createReadStream("ponuda.json");
        var rl = readline.createInterface({
            input: stream,
            crlfDelay: Infinity
        });
        for await (const line of rl) {
            const game = JSON.parse(line);
            games.push(game);
            console.log(game);
        }
    } catch(err) {
        console.error(err);
    }
}

loadJson();

app.get("/games/details/:id", (req, res) => {
    const id = req.params.id;
    for(const game of games) {
        if(game.baseId == id) {
            return res.status(200).json(game);
        }
    }
    return res.status(404).json();
});

app.listen(8000, () => console.log('server up'));
