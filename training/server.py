import sklearn.svm as svm
import sklearn.preprocessing as prep
import pickle
import util
import json

model: svm.OneClassSVM = pickle.load(open("./probe.model", "r"))
scaler: prep.RobustScaler = pickle.load(open("./probe.model", "r"))

def serve():
    sock = util.create_socket()
    sock.bind(util.addr)
    sock.listen()
    while True:
        print("listening for events")
        (connection, addr) = sock.accept()
        print('recieved request from:', addr)
        request = util.read_and_deserialize(connection)
        probeGroup = json.loads(request)
    