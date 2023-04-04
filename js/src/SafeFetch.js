/* 
    Written by Prizzle#4655,
    Method discovered by Milkyway and Sam,
    all rights reserved.
*/

import fetch from "node-fetch";

export default class SafeFetch {
	#maxAttemps = 50;
	#log = (text) => console.log(`[${new Date().toISOString()}] ${text}`);
	#counter = 0;

	async #generateSafeFetchToken(url) {
		try {
			this.#log(`Generating safe fetch token for ${url}`);
			let fetchRaw = await fetch(
				`https://docs.google.com/gview?url=${url}`,
				{
					headers: {
						accept: "*/*",
						"accept-language":
							"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7",
						"sec-ch-ua":
							'" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"',
						"sec-ch-ua-mobile": "?0",
						"sec-ch-ua-platform": '"Windows"',
						"sec-fetch-dest": "empty",
						"sec-fetch-mode": "cors",
						"sec-fetch-site": "same-origin",
						Referer: `https://docs.google.com/gview?url=${url}`,
						"Referrer-Policy": "strict-origin-when-cross-origin",
					},
				}
			);
			if (!fetchRaw.ok || fetchRaw.status == 204) {
				this.#log("Could not generate safe fetch token...");
				return {
					success: false,
					message: `Request was not OK or just empty (${fetchRaw.status})`,
					content: fetchRaw,
				};
			}

			if (fetchRaw.redirected) {
				return {
					message: `Request was unexpectedly redirected to ${fetchRaw.url}`,
					success: false,
					content: fetchRaw,
				};
			}

			let fetchText = await fetchRaw.text();

			//console.log(fetchText);

			if (!fetchText) {
				throw new Error("Fetch got no text");
			}

			let id = fetchText
				?.split(`text?id\\u003d`)?.[1]
				?.split("\\u0026authuser")?.[0]
				?.split('"')[0];

			if (!id) {
				throw new Error("Could not find id");
			}

			return {
				success: true,
				message: "Successfully fetched id",
				content: id,
			};
		} catch (err) {
			return {
				success: false,
				message: `${err.name}, ${err.message}`,
				content: "",
			};
		}
	}
	async #renderSafeFetch(safeFetchId) {
		try {
			let fetchRaw = await fetch(
				`https://docs.google.com/viewerng/text?id=${safeFetchId}&page=0`,
				{
					/* agent: this.#proxyAgent, */
					headers: {
						accept: "*/*",
						"accept-language":
							"de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7",
						"sec-ch-ua":
							'" Not A;Brand";v="99", "Chromium";v="101", "Google Chrome";v="101"',
						"sec-ch-ua-mobile": "?0",
						"sec-ch-ua-platform": '"Windows"',
						"sec-fetch-dest": "empty",
						"sec-fetch-mode": "cors",
						"sec-fetch-site": "same-origin",
					},
				}
			);

			if (!fetchRaw.ok) {
				this.#log("Could not generate safe fetch view...");
				return {
					success: false,
					message: `Request was not OK (${fetchRaw.status})`,
					content: fetchRaw,
				};
			}

			if (fetchRaw.redirected) {
				return {
					message: `Request was unexpectedly redirected to ${fetchRaw.url}`,
					success: false,
					content: fetchRaw,
				};
			}

			let fetchText = await fetchRaw.text();
			let mutateText = fetchText.split(`\n`)[1];
			let fetchJSON = JSON.parse(mutateText);

			if (
				fetchJSON.mimetype == "text/html" &&
				typeof DOMParser != "undefined"
			) {
				let resultHTML = new DOMParser().parseFromString(
					fetchJSON.data,
					"text/html"
				);
				return {
					message:
						"Successfully requested URL and processed content to HTML",
					success: true,
					content: resultHTML,
				};
			} else if (fetchJSON.mimetype == "application/json") {
				try {
					let resultJSON = JSON.parse(fetchJSON.data);
					return {
						message:
							"Successfully requested URL and processed content to JSON",
						success: true,
						content: resultJSON,
					};
				} catch (err) {
					let text = fetchJSON.data;

					return {
						message:
							"Successfully requested URL but failed to parse JSON",
						success: true,
						content: text,
					};
				}
			} else {
				return {
					message: "Successfully requested URL",
					success: true,
					content: fetchJSON,
				};
			}
		} catch (err) {
			return {
				success: false,
				message: `${err.name}, ${err.message}`,
				content: null,
			};
		}
	}
	async get(url) {
		// Attempt to generate a token
		// for the next request
		let tokenGenerationRequest;
		let attempts = 0;
		while (true) {
			tokenGenerationRequest = await this.#generateSafeFetchToken(url);
			attempts++;
			if (tokenGenerationRequest.success) break;
			if (attempts > this.#maxAttemps)
				return {
					message:
						"Exceeded maximum safe fetch token generation attempts",
					success: false,
					content: null,
				};
		}

		// Render your requested content
		// with the generated token
		let renderedContent = await this.#renderSafeFetch(
			tokenGenerationRequest["content"]
		);
		this.#log(
			`[${this.#counter}] Successfully rendered safe fetch content`
		);
		this.#counter++;
		return renderedContent;
	}

	setAttemptLimit(number = this.#maxAttemps) {
		if (!number) return false;
		if (isNaN(parseInt(number))) return false;
		this.#maxAttemps = number;
	}
}
