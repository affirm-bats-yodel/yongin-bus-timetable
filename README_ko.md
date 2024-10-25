# yongin-bus-timetable

Go (Golang) 을 사용한 용인 버스 시간표 추출기.

## 왜 개발하게 되었는가?

* 에버랜드와 수원역을 왔다 갔다 하는 `66` 번, 혹은 `66-4` 번 버스 같은 경우에는
배차간격이 몇십분에 불과해서 조금만 기다리면 타는데 문제가 없음.
* 다만, 다른 배차간격이 넓은 버스를 이용해야 하는 경우에는 정해진 시간표가 있어서
해당 시간표를 모르면 이용이 어려움. (물론 뭐 배차간격과 첫차, 막차로 시간표를 유추할
수도 있긴 한데, 귀찮음.)
* 카카땡, 네땡버 같은 플랫폼에서 제공하는 서비스들은 배차간격이 짧다고 가정하고 시간을 
제공하기 때문에 알기 어렵다는 문제가 있음.
* 찾다 보니, "용인공용버스터미널" 이라는 곳에서 용인시 내에서 운행하고 있는 시내버스를
포함한 다른 버스들의 시간표를 제공함.

## 동작 과정

1. "용인공용버스터미널" 홈페이지의 "시내버스" 에 들어가서, 해당 페이지에서 `button` 컴포넌트
의 `onclick` 속성의 `window.open` 코드의 주소를 추출한다.
2. 추출된 주소로 접속하여 해당 페이지에서 `table` component 의 데이터를 추출한다.
3. 추출된 정보를 구조화 하여 DB (Database, 데이터베이스) 에 집어넣던지, `json` 형태의 데이터로
만들어 CDN 에서 해당 데이터에 접근이 가능하도록 처리한다.
4. PROFIT!