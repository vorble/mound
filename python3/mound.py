import json
import pathlib
import os
import uuid

def uuidToDocDirArray(did):
    return [did[0:2], did[2:4], did[4:6], did[6:8], did]

class Mound:
    MOUND_DATA_DIR = '/tmp/mound'

    @staticmethod
    def setup(MOUND_DATA_DIR):
        Mound.MOUND_DATA_DIR = MOUND_DATA_DIR

    def __init__(self, program, version):
        self.did = uuid.uuid4().hex
        self.program = program
        self.version = version
        self.status = -1
        self.blobs = []
        self.sources = []

        self._open()

    def _open(self):
        docDir = os.path.join(Mound.MOUND_DATA_DIR, *uuidToDocDirArray(self.did))
        os.makedirs(docDir)
        self._writeDoc()

    def close(self, status):
        self.status = status
        self._writeDoc()

    def link(self, sourceDID):
        if sourceDID not in self.sources:
            self.source.append(sourceDID)
        self._writeDoc()

    def _writeDoc(self):
        docPath = os.path.join(Mound.MOUND_DATA_DIR, *uuidToDocDirArray(self.did), 'doc')
        with open(docPath, 'w') as fout:
            fout.write(json.dumps({
                'did': self.did,
                'program': self.program,
                'version': self.version,
                'status': self.status,
                'blobs': self.blobs,
                'sources': self.sources,
            }, separators = (',', ':')))
            fout.write('\n')

    def blob(self, *args):
        if len(args) > 1:
            raise Exception('Max one argument expected.')
        bno = len(self.blobs)
        name = bno if len(args) == 0 else args[0]
        self.blobs.append(name)
        self._writeDoc()
        return bno

    def write(self, bno, data):
        blobPath = os.path.join(Mound.MOUND_DATA_DIR, *uuidToDocDirArray(self.did), str(bno))
        with open(blobPath, 'ab') as fout:
            if isinstance(data, str):
                fout.write(data.encode('utf-8'))
            else:
                fout.write(data)

    def println(self, bno, data):
        self.write(bno, data)
        self.write(bno, '\n')
