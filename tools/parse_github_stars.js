// temp1 is <ul> element with github links on https://github.com/avelino/awesome-go

Promise.all(Array.from(temp1.querySelectorAll('a'))
    .map(e => e.href)
    .filter(link => link.startsWith('https://github.com/'))
    .map(async link => {
        const repo = link.replace('https://github.com/', '')
        const stars = await fetch(`https://api.github.com/repos/${repo}`)
            .then(r => r.json())
            .then(d => d.stargazers_count)
        return {repo, stars}
    }))
    .then(result => result.sort((a, b) => b.stars - a.stars))
    .then(console.log)
