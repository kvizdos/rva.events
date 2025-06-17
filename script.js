const title = document.querySelector("h1").innerText;

let output = [];

const dateTime = document.querySelector("time.block").innerText;
const [date, time] = dateTime.split(" • ");

let venue = document.querySelector("[data-testid='venue-name-link']");
let loc = document.querySelector("[data-testid='location-info']");

let summaryParts = [];
document.querySelectorAll("#event-details p").forEach((el) => {
	summaryParts.push(el.innerText);
});

const summary = summaryParts.join(" ").replace(/\s+/g, " ").trim();

let metadata = [
	"---",
	`Title: ${title}`,
	`Date: ${date.split("\n")[0]}`,
	`Time: ${date.split("\n")[1]}`,
	`Summary: "${summary}"`,
	`Author: unknown`, // fill in if you can scrape
	`Tags: Tech, Free`, // adjust as needed
	`Price: FREE`, // adjust as needed
	`Location: ${venue?.innerText || "Unknown"} (${loc?.innerText || "Unknown"})`,
	`RSVP: ${location.href}`,
	"---",
];

output.push("## WHY GO", ...summaryParts);
output.push("## WHEN TO GO", dateTime);
output.push("## WHERE TO GO", venue?.innerText || "", loc?.innerText || "");

console.log([...metadata, "", ...output].join("\n\n"));
