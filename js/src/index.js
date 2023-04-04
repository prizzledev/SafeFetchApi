import express from "express";
const app = express();
import SafeFetch from "./SafeFetch.js";
import cors from "cors";

app.use(cors());
app.use(express.json());

function log(message) {
	console.log(`[${new Date().toISOString()}] ${message}`);
}

let safeFetch = new SafeFetch();

app.post("/api/v2/safeFetch", async (req, res) => {
	log("Request received");

	let url;

	try {
		url = req.body.url;
	} catch (err) {
		res.status(400).send({
			message: "Invalid request body",
			success: false,
			content: null,
		});
		return;
	}

	let result;

	// result = key
	result = await safeFetch.get(url);

	if (!result.success) {
		res.status(500).send({
			message: result.message,
			success: false,
			content: null,
		});
		return;
	}

	res.send({
		message: "Successfully requested URL",
		success: true,
		content: result.content,
	});
});

app.listen(4501, () => {
	log("server is listening http://localhost:4051/");
});
