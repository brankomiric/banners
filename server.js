const express = require('express');
const cors = require("cors");
const fs = require('fs');
const readline = require('readline');

const app = express();
app.use(cors());
app.use(express.json());

const gamesDetails = [];
let games = {};

const loadGameDetails = async () => {
    try {
        const stream = fs.createReadStream("ponuda.json");
        var rl = readline.createInterface({
            input: stream,
            crlfDelay: Infinity
        });
        for await (const line of rl) {
            const game = JSON.parse(line);
            gamesDetails.push(game);
            console.log(game);
        }
    } catch(err) {
        console.error(err);
    }
}

const loadGames = () => {
    let content = fs.readFileSync("promo_banners.json");
    content = JSON.parse(content);
    games = content;
}

loadGameDetails();
loadGames();

app.get("/promo_banners", (req, res) => {
    return res.status(200).json(games);
});

app.post("/games/details", (req, res) => {
    const ids = req.body.matchIds;
    const matchedGameDetails = [];
    for(id of ids) {
        for(const game of gamesDetails) {
            if(game.baseId == id) {
                matchedGameDetails.push(game);
                break;
            }
        }
    }
    return res.status(200).json(matchedGameDetails);
});

app.listen(8000, () => console.log('server up'));
