from python3.mound import Mound

Mound.setup('/home/keith/mound_data')

m = Mound('mound', '0.0.1-test')
b0 = m.blob()
m.println(b0, 'Hello, python!')
b1 = m.blob()
m.println(b1, 'Goodbye, python!')
m.close(0)
