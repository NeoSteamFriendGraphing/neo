import { setCrawlPageUserCardDetails } from '/static/javascript/userCard.js';

export function doesProcessedGraphDataExist(crawlID) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2590/api/doesprocessedgraphdataexist/${crawlID}`, {
            method: 'POST',
            headers: {
                "Content-Type": "application/json"
            },
        }).then((res => res.json()))
        .then(data => {
            if (data.exists == "yes") {
                resolve(true)
            } 
            resolve(false)
        }).catch(err => {
            reject(err)
        })
    });
}

export function getProcessedGraphData(crawlID) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2590/api/getprocessedgraphdata/${crawlID}`, {
            method: "POST",
            headers: {
                'Accept-Encoding': 'gzip'
            }
        }).then(res => res.json())
        .then(data => {
            resolve(data)
        }).catch(err => {
            reject(err)
        })
    });
}

export function getAnyNewFinishedCrawlStatuses() {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2590/api/getfinishedcrawlsaftertimestamp?timestamp=${5}`, {
            headers: {
                "Content-Type": "application/json"
            },
        }).then((res => res.json()))
        .then(data => {
            resolve(data.crawls)
        }).catch(err => {
            reject(err)
        })
    })
}

export function getAnyNewFinishedShortestDistanceCrawlStatuses() {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2590/api/getfinishedshortestdistancecrawlsaftertimestamp?timestamp=${-5}`, {
            headers: {
                "Content-Type": "application/json"
            },
        }).then((res => res.json()))
        .then(data => {
            resolve(data.crawlingstatus)
        }).catch(err => {
            reject(err)
        })
    })
}

export function getUser(crawlID) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2590/api/getcrawlinguser/${crawlID}`, {
            headers: {
                "Content-Type": "application/json"
            },
        }).then((res => res.json()))
        .then(data => {
            if (data.error) {
                reject(data)
            }
            resolve(data.user)
        }).catch(err => {
            console.error(err);
        })
    })
}

export function getShortestDistanceInfo(crawlIDs) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2590/api/getshortestdistanceinfo`, {
            method: 'POST',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({"crawlids":crawlIDs})
        }).then((res => res.json()))
        .then(data => {
            if (data.error) {
                reject(data)
            }
            resolve(data.shortestdistanceinfo)
        }).catch(err => {
            console.error(err);
        })
    })
}

export function isCrawlingFinished(crawlID) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2590/api/getcrawlingstatus/${crawlID}`, {
            headers: {
                "Content-Type": "application/json"
            }
        }).then((res => res.json()))
        .then(data => {
            if (data.crawlingstatus.totaluserstocrawl == data.crawlingstatus.userscrawled) {
                resolve(true)
            }
            resolve(false)
        }).catch(err => {
            reject(err)
        })
    })
}

export function getCrawlingStatus(crawlID) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2590/api/getcrawlingstatus/${crawlID}`, {
            headers: {
                "Content-Type": "application/json"
            }
        }).then((res => res.json()))
        .then(data => {
            resolve(data.crawlingstatus)
        }).catch(err => {
            reject(err)
        })
    })
}


export function startCreateGraph(crawlID) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2570/creategraph/${crawlID}`, {
            method: 'POST',
            headers: {
                "Content-Type": "application/json"
            },
        }).then((res => res.json()))
        .then(data => {
            resolve(data)
        }).catch(err => {
            console.error(err);
        })
    })
}

export function waitUntilGraphDataExists(crawlID) {
    return new Promise((resolve, reject) => {
        let interval = setInterval(function() {
            doesProcessedGraphDataExist(crawlID).then(doesExist => {
                if (doesExist === true) {
                    clearInterval(interval);
                    resolve(true)
                }
            }, err => {
                clearInterval(interval);
                reject(err)
            })
        }, 1000);
    })
}

export function getCrawlingUserWhenAvailable(crawlID , idPrefix) {
    return new Promise((resolve, reject) => {
    let interval = setInterval(() => {
        if (usersCrawledIsMoreThanZero(idPrefix)) {
            clearInterval(interval);
            getUser(crawlID).then(user => {
                setCrawlPageUserCardDetails(user, idPrefix)
                resolve(true)
            }, err => {
                reject(err)
            })
        }
      }, 50);
    });
}

export function startCalculateGetShortestDistance(crawlIDs) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2590/api/calculateshortestdistanceinfo`, {
            method: 'POST',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({"crawlids": crawlIDs})
        }).then((res => res.json()))
        .then(data => {
            resolve(data)
        }).catch(err => {
            reject(err)
        })
    })
}

export function startCrawl(crawlDTO) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2570/crawl`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(crawlDTO),
        }).then(res => res.json())
        .then(data => {
            resolve(data.crawlids)
        }).catch(err => {
            reject(err)
        })
    })
}

export function hasBeenCrawled(steamID, level) {
    return new Promise((resolve, reject) => {
        const reqBody = {
            "level": parseInt(level),
            "steamid": steamID
        }
        fetch(`http://localhost:2590/api/hasbeencrawledbefore`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(reqBody),
        }).then(res => res.json())
        .then(data => {
            resolve(data.message)
        }).catch(err => {
            reject(err)
        })
    })
} 

export function isPublicProfile(steamID) {
    return new Promise((resolve, reject) => {
        fetch(`http://localhost:2570/isprivateprofile/${steamID}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
    }).then(res => res.json())
    .then(data => {
        if (data.message === "public") {
            resolve(true)
        }
        resolve(false)
    }).catch(err => {
        console.error(err)
        reject(err)
        })
    })
}

function usersCrawledIsMoreThanZero(idPrefix) {
    if (parseInt(document.getElementById(`${idPrefix}UsersCrawled`).textContent) >= 1) {
        return true
    }
    return false
}