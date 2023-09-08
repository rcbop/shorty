document.addEventListener("DOMContentLoaded", function () {
    var form = document.getElementById("urlForm");
    var inputField = document.getElementById("longURL");
    var shortenedURL = document.getElementById("shortenedURL");
    var shortURLDisplay = document.getElementById("shortURL");
    var copyButton = document.getElementById("copyButton");

    form.addEventListener("submit", function (event) {
        event.preventDefault();

        var longURL = inputField.value;
        console.log(">> POST longURL: " + longURL);

        if (!isValidURL(longURL)) {
            alert("Please enter a valid URL (e.g., http://example.com)");
            return;
        }

        fetch("/shorty", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ longURL: longURL }),
        })
            .then(function (response) {
                if (!response.ok) {
                    throw new Error("Network response was not ok");
                }
                return response.json();
            })
            .then(function (data) {
                var message = data.message;
                console.log("<< Response: " + message);

                form.style.display = "none";
                shortenedURL.style.display = "block";
                shortURLDisplay.textContent = message;
                copyButton.style.display = "block";
            })
            .catch(function (error) {
                console.error("Fetch Error:", error.message);
            });
    });

    copyButton.addEventListener("click", function () {
        var textArea = document.createElement("textarea");
        textArea.value = shortURLDisplay.textContent;
        document.body.appendChild(textArea);
        textArea.select();
        document.execCommand("copy");
        document.body.removeChild(textArea);
        alert("Shortened URL copied to clipboard!");
    });

    function isValidURL(url) {
        var pattern = /^(http|https):\/\/[^ "]+$/;
        return pattern.test(url);
    }
});
