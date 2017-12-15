import sys
import numpy as np


def main():
    vertices = []
    neighbours = []
    with open(sys.argv[1]) as f:
        for line in f:
            tokens = line.split()
            if len(tokens) == 0:
                continue

            kind = tokens[0]
            if kind == "v":
                v = np.array(tokens[1:]).astype(float)
                vertices.append(v)
                neighbours.append([])
            elif kind == "f":
                f = np.array(tokens[1:]).astype(int) - 1
                for i, v in enumerate(f):
                    l = f[i - 1]
                    r = f[(i + 1) % len(f)]
                    neighbours[v].append((vertices[l], vertices[r]))

    with open(sys.argv[1], 'a') as f:
        norms = []
        for v, n in zip(vertices, neighbours):
            s = []
            for (l, r) in n:
                n = np.cross(r - v, l - v)
                n /= np.linalg.norm(n)
                vl = (l - v) / np.linalg.norm(l - v)
                vr = (r-v) / np.linalg.norm(r - v)
                s.append(n * np.arccos(np.dot(vl, vr)) / 2 / np.pi)
            s = np.average(s, axis=0)
            s /= np.linalg.norm(s)
            norms.append(s)

        f.writelines(list(map(lambda n: "vn {} {} {}\n".format(n[0], n[1], n[2]), norms)))


if __name__ == '__main__':
    main()
