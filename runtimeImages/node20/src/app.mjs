import express from "express";
import handler from './handler.mjs';
const app = express();

app.get('/*', (req, res) => {
    handler(req, res);
});

app.listen(443, () => {
    console.log(`Server listening at http://localhost:443`);
});
