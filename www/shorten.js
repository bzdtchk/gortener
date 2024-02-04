document.getElementById('link').addEventListener('click', async function (event) {
    let element = document.getElementById('link');
    try {
        await navigator.clipboard.writeText(element.innerText);
        alert('Text byl zkopírován');
    } catch (err) {
        console.error('Failed to copy: ', err);
    }
});

document.getElementById('form').addEventListener('submit', function (event) {
    event.preventDefault();
    document.getElementById('result-error').style.display = 'none';
    document.getElementById('result-success').style.display = 'none';

    let destinationUrlValue = document.getElementsByName('destination_url')[0].value;
    let slugValue = document.getElementsByName('slug')[0].value;

    fetch("api/shorten-link", {
        method: "POST",
        body: JSON.stringify({
            destination_url: destinationUrlValue,
            slug: slugValue
        }),
        headers: {
            "Content-type": "application/json; charset=UTF-8"
        }
    }).then(async response => {
        if (!response.ok) {
            let responseJson = await response.json();
            throw new Error(responseJson?.message ?? 'Chyba!');
        }
        return response.json();
    }).then(data => {
        document.getElementById('result-success').style.display = 'initial';
        document.getElementById('link').innerText = window.location.protocol + "//" + window.location.host + "/" + data.slug;
        document.getElementById('form').reset();
    }).catch(error => {
        document.getElementById('result-error').style.display = 'initial';
        document.getElementById('result-error').innerHTML = parseMessage(error.message);
    });
});

function parseMessage(message) {
    let messages = {
        'Invalid JSON received': 'Vyplňte odkaz ke zkrácení i vlastní zkratku',
        'Destination URL for shortening is not reachable': 'Odkaz ke zkrácení vede na stránku, která není dostupná.<br>Zkontrolujte správnost odkazu',
        'There\'s already link with the same slug in our database': 'Tato vlastní zkratka je již zabraná, prosím zadejte jinou',
        'Slug must be an alphanumerical string': 'Vlastní zkratka musí být alfanumerická (písmenka bez diakritiky s čísly, spojené pomlčkou). Například toto-je-spravna-zkratka'
    };

    return messages[message] ?? message;
}