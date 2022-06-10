const MAX_RESULTS_PER_ROW = 12;

function presentResults(results) {
    let body = document.getElementById("anagram_results");
    body.innerHTML = "";
    const ordered = Object.keys(results).sort().reduce(
        (obj, key) => {
            obj[key] = results[key];
            return obj;
        },
        {}
    );
    for (let key in ordered) {
        let table = document.createElement("table");
        let thead = document.createElement("thead");
        let tbody = document.createElement("tbody");
        table.appendChild(thead);
        table.appendChild(tbody);
        body.appendChild(table);
        let caption = table.createCaption();
        caption.innerHTML = `<strong>${key}-character matches:</strong>`
        let words = ordered[key]
        let target_rows = Math.ceil(words.length / MAX_RESULTS_PER_ROW);
        for (let i = 0; i < target_rows; i++) {
            let row = table.insertRow(-1);
            tmp = words.splice(0, MAX_RESULTS_PER_ROW);
            for (word_index in tmp) {
                let cell = row.insertCell(word_index);
                cell.innerHTML = tmp[word_index];
            }
        }
    }
}

const anagramForm = document.getElementById("anagram_input");
const anagramSubmitButton = document.getElementById("anagram_submit");

anagramSubmitButton.addEventListener("click", (e) => {
    console.time("anagramSubmitButton")
    e.preventDefault();
    const query_word = anagramForm.anagram_input.value;
    const partial = anagramForm.include_partial.checked;
    let form = new FormData();
    form.append("word", query_word);
    form.append("partial", partial);
    (async () => {
        const resp = await fetch("/anagrams", {
            body: form,
            method: "POST"
        });
        const content = await resp.json();
        presentResults(content);
    })()
    console.timeEnd("anagramSubmitButton")
})