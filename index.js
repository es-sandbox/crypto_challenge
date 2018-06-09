var xhr = new XMLHttpRequest();

xhr.open('GET', 'http:/localhost:9000/order/sell', true);

xhr.send(); // (1)

xhr.onreadystatechange = function() { // (3)
    if (xhr.readyState != 4) return;

    if (xhr.status != 200) {
        console.log(xhr.status + ': ' + xhr.statusText);
    } else {
        // console.log(xhr.responseText);

        var table = document.getElementById("sell-order");
        var sellOrderList = JSON.parse(xhr.responseText);
        sellOrderList.forEach(function(item, i, sellOrderList) {
            if (i >= 20) {
                return
            }

            table.innerHTML += "<tr><td>" + convertOrderTypeToString(item.OrderType) + "</td><td>" + item.Price + "</td><td>" + item.Amount + "</td></tr>"
        });
    }
};

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var xhr2 = new XMLHttpRequest();

xhr2.open('GET', 'http:/localhost:9000/order/buy', true);

xhr2.send(); // (1)

xhr2.onreadystatechange = function() { // (3)
    if (xhr2.readyState != 4) return;

    if (xhr2.status != 200) {
        console.log(xhr2.status + ': ' + xhr2.statusText);
    } else {
        // console.log(xhr2.responseText);

        var table = document.getElementById("buy-order");
        var buyOrderList = JSON.parse(xhr2.responseText);
        buyOrderList.forEach(function(item, i, buyOrderList) {
            if (i >= 20) {
                return
            }

            table.innerHTML += "<tr><td>" + convertOrderTypeToString(item.OrderType) + "</td><td>" + item.Price + "</td><td>" + item.Amount + "</td></tr>"
        });
    }
};

function convertOrderTypeToString(orderType) {
    if (orderType == 0) {
        return "Sell"
    }
    return "Buy"
}

var xhr3 = new XMLHttpRequest();

xhr3.open('GET', 'http:/localhost:9000/deals', true);

xhr3.send(); // (1)

xhr3.onreadystatechange = function() { // (3)
    if (xhr3.readyState != 4) return;

    if (xhr3.status != 200) {
        console.log(xhr3.status + ': ' + xhr3.statusText);
    } else {
        // console.log(xhr3.responseText);

        var table = document.getElementById("deals");
        var deals = JSON.parse(xhr3.responseText);
        deals.forEach(function(item, i, deals) {
            if (i + 20 < deals.length) {
                return
            }

            console.log(item);
            table.innerHTML += "<tr><td>" + item.AmountBTC + "</td><td>" + item.Price + "</td></tr>"
        });
    }
};

var xhr4 = new XMLHttpRequest();

xhr4.open('GET', 'http:/localhost:9000/price', true);

xhr4.send(); // (1)

xhr4.onreadystatechange = function() { // (3)
    if (xhr4.readyState != 4) return;

    if (xhr4.status != 200) {
        console.log(xhr4.status + ': ' + xhr4.statusText);
    } else {
        // console.log(xhr3.responseText);

        var table = document.getElementById("price");
        var price = JSON.parse(xhr4.responseText);
        table.innerHTML = price
    }
};