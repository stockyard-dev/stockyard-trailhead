package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Trailhead</title>
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#c45d2c;--rl:#e8753a;--leather:#a0845c;--ll:#c4a87a;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c44040;--mono:'JetBrains Mono',Consolas,monospace;--serif:'Libre Baskerville',Georgia,serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);font-size:13px;line-height:1.6}
.hdr{padding:.6rem 1.2rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}
.hdr h1{font-family:var(--serif);font-size:1rem}.hdr h1 span{color:var(--rl)}
.main{max-width:600px;margin:0 auto;padding:1rem 1.2rem}
.btn{font-family:var(--mono);font-size:.68rem;padding:.3rem .6rem;border:1px solid;cursor:pointer;background:transparent;transition:.15s;white-space:nowrap}
.btn-p{border-color:var(--rust);color:var(--rl)}.btn-p:hover{background:var(--rust);color:var(--cream)}
.btn-d{border-color:var(--bg3);color:var(--cm)}.btn-d:hover{border-color:var(--red);color:var(--red)}
.progress{margin-bottom:1rem;text-align:center}
.progress-bar{height:6px;background:var(--bg3);border-radius:3px;overflow:hidden;margin-top:.3rem}
.progress-fill{height:100%;background:var(--green);transition:width .3s}
.progress-text{font-size:.8rem;color:var(--cream)}
.progress-sub{font-size:.65rem;color:var(--cm)}

.habit{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;margin-bottom:.5rem;display:flex;align-items:center;gap:.7rem;transition:.1s}
.habit:hover{background:var(--bg3)}
.check-btn{width:28px;height:28px;border-radius:50%;border:2px solid var(--bg3);cursor:pointer;display:flex;align-items:center;justify-content:center;transition:.15s;flex-shrink:0;font-size:.8rem}
.check-btn.done{border-color:var(--green);background:var(--green);color:var(--bg)}
.check-btn:hover{border-color:var(--green)}
.habit-info{flex:1}
.habit-name{font-size:.8rem;font-weight:600;display:flex;align-items:center;gap:.3rem}
.habit-meta{font-size:.65rem;color:var(--cm);display:flex;gap:.7rem}
.streak{font-weight:600}.streak.active{color:var(--gold)}
.color-dot{width:8px;height:8px;border-radius:50%;flex-shrink:0}
.habit-actions{display:flex;gap:.3rem}

.modal-bg{position:fixed;top:0;left:0;right:0;bottom:0;background:rgba(0,0,0,.65);display:flex;align-items:center;justify-content:center;z-index:100}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:90%;max-width:400px}
.modal h2{font-family:var(--serif);font-size:.9rem;margin-bottom:.8rem}
label.fl{display:block;font-size:.65rem;color:var(--leather);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem;margin-top:.5rem}
input[type=text],input[type=color],select{background:var(--bg);border:1px solid var(--bg3);color:var(--cream);padding:.35rem .5rem;font-family:var(--mono);font-size:.78rem;width:100%;outline:none}
.empty{text-align:center;padding:2rem;color:var(--cm);font-style:italic;font-family:var(--serif)}
</style>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital@0;1&family=JetBrains+Mono:wght@400;600&display=swap" rel="stylesheet">
</head><body>
<div class="hdr"><h1><span>Trailhead</span></h1><button class="btn btn-p" onclick="showNewHabit()">+ Habit</button></div>
<div class="main">
<div class="progress" id="progress"></div>
<div id="habitList"></div>
</div>
<div id="modal"></div>
<script>
let habits=[];
async function api(url,opts){return(await fetch(url,opts)).json()}
function esc(s){return String(s||'').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;')}

async function load(){
  const d=await api('/api/today');
  habits=d.habits||[];
  const pct=d.total?Math.round(d.done/d.total*100):0;
  document.getElementById('progress').innerHTML=
    '<div class="progress-text">'+d.done+'/'+d.total+' completed today</div>'+
    '<div class="progress-bar"><div class="progress-fill" style="width:'+pct+'%"></div></div>'+
    '<div class="progress-sub">'+d.date+'</div>';
  renderHabits();
}

function renderHabits(){
  const el=document.getElementById('habitList');
  if(!habits.length){el.innerHTML='<div class="empty">No habits yet. Add one to start tracking.</div>';return}
  el.innerHTML=habits.map(h=>{
    const done=h.checked_today;
    return '<div class="habit">'+
      '<div class="check-btn '+(done?'done':'')+'" onclick="toggle(\''+h.id+'\','+done+')">'+(done?'✓':'')+'</div>'+
      '<div class="color-dot" style="background:'+esc(h.color)+'"></div>'+
      '<div class="habit-info">'+
        '<div class="habit-name">'+esc(h.name)+'</div>'+
        '<div class="habit-meta">'+
          '<span class="streak '+(h.streak>0?'active':'')+'">🔥 '+h.streak+'d streak</span>'+
          '<span>Best: '+h.best_streak+'d</span>'+
          '<span>'+h.total_checks+' total</span>'+
        '</div>'+
      '</div>'+
      '<div class="habit-actions">'+
        '<span style="cursor:pointer;font-size:.6rem;color:var(--cm)" onclick="editHabit(\''+h.id+'\')">edit</span>'+
        '<span style="cursor:pointer;font-size:.6rem;color:var(--cm)" onclick="if(confirm(\'Delete?\'))delHabit(\''+h.id+'\')">del</span>'+
      '</div></div>'
  }).join('')
}

async function toggle(id,done){
  if(done){await api('/api/habits/'+id+'/uncheck',{method:'POST',headers:{'Content-Type':'application/json'},body:'{}'})}
  else{await api('/api/habits/'+id+'/check',{method:'POST',headers:{'Content-Type':'application/json'},body:'{}'})}
  load()
}

function showNewHabit(){
  document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal">'+
    '<h2>New Habit</h2><label class="fl">Name</label><input type="text" id="nh-name" placeholder="Exercise">'+
    '<label class="fl">Description</label><input type="text" id="nh-desc">'+
    '<div style="display:flex;gap:.5rem"><div style="flex:1"><label class="fl">Frequency</label><select id="nh-freq"><option value="daily">Daily</option><option value="weekly">Weekly</option></select></div>'+
    '<div><label class="fl">Color</label><input type="color" id="nh-color" value="#c45d2c" style="height:30px"></div></div>'+
    '<div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveNewHabit()">Create</button><button class="btn btn-d" onclick="closeModal()">Cancel</button></div></div></div>';
}
async function saveNewHabit(){
  const body={name:document.getElementById('nh-name').value,description:document.getElementById('nh-desc').value,frequency:document.getElementById('nh-freq').value,color:document.getElementById('nh-color').value};
  if(!body.name){alert('Name required');return}
  await api('/api/habits',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});closeModal();load()
}
function editHabit(id){
  const h=habits.find(x=>x.id===id);if(!h)return;
  document.getElementById('modal').innerHTML='<div class="modal-bg" onclick="if(event.target===this)closeModal()"><div class="modal">'+
    '<h2>Edit Habit</h2><label class="fl">Name</label><input type="text" id="eh-name" value="'+esc(h.name)+'">'+
    '<label class="fl">Description</label><input type="text" id="eh-desc" value="'+esc(h.description)+'">'+
    '<div style="display:flex;gap:.5rem"><div style="flex:1"><label class="fl">Frequency</label><select id="eh-freq"><option value="daily"'+(h.frequency==='daily'?' selected':'')+'>Daily</option><option value="weekly"'+(h.frequency==='weekly'?' selected':'')+'>Weekly</option></select></div>'+
    '<div><label class="fl">Color</label><input type="color" id="eh-color" value="'+esc(h.color)+'" style="height:30px"></div></div>'+
    '<div style="display:flex;gap:.5rem;margin-top:1rem"><button class="btn btn-p" onclick="saveEditHabit(\''+id+'\')">Save</button><button class="btn btn-d" onclick="closeModal()">Cancel</button></div></div></div>';
}
async function saveEditHabit(id){
  const body={name:document.getElementById('eh-name').value,description:document.getElementById('eh-desc').value,frequency:document.getElementById('eh-freq').value,color:document.getElementById('eh-color').value};
  await api('/api/habits/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});closeModal();load()
}
async function delHabit(id){await api('/api/habits/'+id,{method:'DELETE'});load()}
function closeModal(){document.getElementById('modal').innerHTML=''}
load()
</script></body></html>`
