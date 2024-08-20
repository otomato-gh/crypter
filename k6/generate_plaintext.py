import json
import random
import string

def generate_random_text(length):
    letters = string.ascii_letters
    return ''.join(random.choice(letters) for i in range(length))

def main():
    array = []
    for i in range(1,10):
        random_text = generate_random_text(10000)
        data = {"plaintext": random_text}
#        print(data)
        array.append(data)
#        print(array)

    with open("random_text.json", "w") as f:
        json.dump(array, f, indent=4)

if __name__ == "__main__":
    main()
