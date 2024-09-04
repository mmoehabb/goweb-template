gameframe = null
gamelabel = null

ctx = null
started = false
gameloop = null

function start() {
  if (started) return;
  started = true;
  gameframe = document.getElementById("gameframe")
  gamelabel = document.getElementById("gamelabel")
  gamelabel.innerText = "Press anywhere on the canvas to jump"
  ctx = gameframe.getContext("2d")

  var audio = new Audio('public/sounds/bg.mp3');
  audio.loop = true;
  audio.play();

  const obstacles = []
  obstacles.push(getObj('public/PNGs/php.png'))
  obstacles.push(getObj('public/PNGs/js.png'))
  obstacles.push(getObj('public/PNGs/python.png'))
  obstacles.push(getObj('public/PNGs/java.png'))

  gameloop = setInterval(() => {
    ctx.fillRect(0, 0, gameframe.width, gameframe.height)
    for (let obs of obstacles) {
      drawObstacle(ctx, obs)
    }
    drawPlayer(ctx)
  }, 1000/60);
}


const player = {
  position: {
    x: 0,
    y: 0,
  },
  size: {
    width: 0,
    height: 0,
  },
  image: null
}

let dt = 0
let jumpforce = 0
let onground = true
function drawPlayer(ctx) {
  if (!player.image) {
    player.image = new Image()
    player.image.onload = function() {
      player.size.width = gameframe.width / 5
      player.size.height = 0.5 * player.size.width * this.height / this.width
      player.position.x = gameframe.width / 5
      player.position.y = gameframe.height - player.size.height
    }
    player.image.src = "public/PNGs/player.png"
    window.onmousedown = () => {
      if (onground) {
        jumpforce = 50
        onground = false
      }
    }
  }

  // gravity
  player.position.y = gameframe.height - player.size.height - (jumpforce * dt)
  if (jumpforce > 0) {
    dt += 0.05
    jumpforce -= dt*dt
  }
  else {
    dt = 0
    jumpforce = 0
    onground = true
  }

  ctx.drawImage(
    player.image, 
    player.position.x,
    player.position.y, 
    player.size.width, 
    player.size.height 
  )
}

function getObj(imgpath) {
  const obj = { position: {}, size: {} }
  if (!obj.image) {
    obj.image = new Image()
    obj.image.onload = function() {
      obj.size.height = gameframe.height / 7
      obj.size.width = 2 * obj.size.height * this.width / this.height 
      obj.position.x = -(obj.size.width + 20)
      obj.position.y = gameframe.height - obj.size.height
    }
    obj.image.src = imgpath
    obj.speed = 1
  }
  return obj
}

function drawObstacle(ctx, obj) {
  // movement
  obj.position.x -= obj.speed
  if (obj.position.x < -obj.size.width) {
    obj.position.x = gameframe.width * (Math.floor(Math.random() * 10) % 9 + 2)
    obj.speed = Math.floor(Math.random() * 10) % 3 + 3
  }

  // trigger collision with player
  if (triggerCollision(player, obj)) {
    gameover()
  }

  ctx.drawImage(
    obj.image, 
    obj.position.x,
    obj.position.y, 
    obj.size.width, 
    obj.size.height 
  )
}

function gameover() {
  clearInterval(gameloop)
  gamelabel.innerText = "Oh... I stepped in shit"
}

function triggerCollision(obj1, obj2) {
  const sobj = obj1.width <= obj2.width ? obj1 : obj2
  const bobj = obj1.width <= obj2.width ? obj2 : obj1
  // corners of small obj
  const corners = [
    {
      x: sobj.position.x,
      y: sobj.position.y,
    },
    {
      x: sobj.position.x + sobj.size.width,
      y: sobj.position.y,
    },
    {
      x: sobj.position.x,
      y: sobj.position.y + sobj.size.height,
    },
    {
      x: sobj.position.x + sobj.size.width,
      y: sobj.position.y + sobj.size.height,
    },
  ]
  const inpadding = 10
  for (let corner of corners) {
    if (corner.x >= bobj.position.x + inpadding && corner.x <= bobj.position.x + bobj.size.width - inpadding) {
      if (corner.y >= bobj.position.y + inpadding && corner.y <= bobj.position.y + bobj.size.height - inpadding) {
        return true
      }
    }
  }
  return false
}
