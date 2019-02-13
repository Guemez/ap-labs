def substring(word):
    count = 0
    max_sub = 0
    seen = {}
    for i in range(len(word)):
        if word[i] not in seen:
            seen[word[i]] = i
            count += 1
        else:
            if max_sub < count:
                max_sub = count
            tmp = list(seen.keys())
            if word[i] != tmp[0]:
                del seen[tmp[0]]
                count -= 1
            del seen[word[i]]
            seen[word[i]] = i
    return max(count, max_sub)

print(substring("abbb"))
