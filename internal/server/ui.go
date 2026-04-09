package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Trailhead</title>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.6;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-family:var(--serif);font-size:1rem;letter-spacing:1px}
.hdr h1 span{color:var(--rust)}
.main{max-width:680px;margin:0 auto;padding:1.2rem 1.5rem}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1.2rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.3rem;font-weight:700;color:var(--gold)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.progress{margin-bottom:1rem;background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem}
.progress-text{font-size:.85rem;color:var(--cream);font-weight:700;display:flex;justify-content:space-between;align-items:center}
.progress-pct{color:var(--gold)}
.progress-bar{height:6px;background:var(--bg3);overflow:hidden;margin-top:.4rem}
.progress-fill{height:100%;background:var(--green);transition:width .3s}
.progress-sub{font-size:.55rem;color:var(--cm);margin-top:.3rem;text-transform:uppercase;letter-spacing:1px}
.toolbar{display:flex;gap:.5rem;margin-bottom:.8rem;align-items:center}
.toolbar .arch-toggle{font-size:.55rem;color:var(--cm);display:inline-flex;align-items:center;gap:.3rem;cursor:pointer;user-select:none}
.toolbar .arch-toggle input{width:auto;cursor:pointer}
.habit{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem .8rem;margin-bottom:.5rem;display:flex;align-items:center;gap:.7rem;transition:.1s}
.habit:hover{border-color:var(--leather)}
.habit.archived{opacity:.5}
.check-btn{width:32px;height:32px;border-radius:50%;border:2px solid var(--bg3);cursor:pointer;display:flex;align-items:center;justify-content:center;transition:.15s;flex-shrink:0;font-size:.9rem;background:transparent;color:var(--bg)}
.check-btn:hover{border-color:var(--green)}
.check-btn.done{border-color:var(--green);background:var(--green);color:var(--bg)}
.color-dot{width:8px;height:8px;border-radius:50%;flex-shrink:0}
.habit-info{flex:1;min-width:0}
.habit-name{font-size:.82rem;font-weight:700;display:flex;align-items:center;gap:.4rem}
.habit-desc{font-size:.65rem;color:var(--cm);font-style:italic;margin-top:.1rem;line-height:1.4}
.habit-meta{font-size:.6rem;color:var(--cm);display:flex;gap:.7rem;flex-wrap:wrap;margin-top:.25rem;align-items:center}
.streak{font-weight:700}
.streak.active{color:var(--gold)}
.habit-extra{font-size:.55rem;color:var(--cd);margin-top:.4rem;padding-top:.3rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.habit-extra-row{display:flex;gap:.4rem}
.habit-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.habit-extra-val{color:var(--cream)}
.habit-actions{display:flex;gap:.4rem;flex-shrink:0}
.icon-btn{font-size:.55rem;color:var(--cm);cursor:pointer;background:none;border:none;font-family:var(--mono);padding:.2rem .3rem;text-transform:uppercase;letter-spacing:.5px}
.icon-btn:hover{color:var(--cream)}
.icon-btn.del:hover{color:var(--red)}
.btn{font-family:var(--mono);font-size:.65rem;padding:.3rem .65rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-family:var(--serif);font-size:.95rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.fr input[type=color]{height:2rem;padding:.2rem}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto;color:var(--red);border-color:#3a1a1a}
.acts .btn-del:hover{border-color:var(--red);color:var(--red)}
.empty{text-align:center;padding:2rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}.trial-bar{flex-direction:column;align-items:stretch}.trial-bar input.key-input{width:100%}}
.trial-bar{display:none;background:linear-gradient(90deg,#3a2419,#2e1c14);border-bottom:2px solid var(--rust);padding:.7rem 1.5rem;font-family:var(--mono);font-size:.68rem;color:var(--cream);align-items:center;gap:1rem;flex-wrap:wrap}
.trial-bar.show{display:flex}
.trial-bar-msg{flex:1;min-width:240px;line-height:1.5}
.trial-bar-msg strong{color:var(--rust);text-transform:uppercase;letter-spacing:1px;font-size:.6rem;display:block;margin-bottom:.15rem}
.trial-bar-actions{display:flex;gap:.5rem;align-items:center;flex-wrap:wrap}
.trial-bar a.btn-trial{background:var(--rust);color:#fff;padding:.4rem .8rem;text-decoration:none;font-size:.65rem;text-transform:uppercase;letter-spacing:1px;font-weight:700;border:1px solid var(--rust);transition:all .2s}
.trial-bar a.btn-trial:hover{background:#f08545;border-color:#f08545}
.trial-bar-divider{color:var(--cm);font-size:.6rem}
.trial-bar input.key-input{padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.6rem;width:200px}
.trial-bar input.key-input:focus{outline:none;border-color:var(--rust)}
.trial-bar button.btn-activate{padding:.4rem .7rem;background:var(--bg2);color:var(--cream);border:1px solid var(--leather);font-family:var(--mono);font-size:.6rem;cursor:pointer;text-transform:uppercase;letter-spacing:1px}
.trial-bar button.btn-activate:hover{background:var(--bg3)}
.trial-bar button.btn-activate:disabled{opacity:.5;cursor:wait}
.trial-msg{font-size:.6rem;color:var(--cm);margin-left:.5rem}
.trial-msg.error{color:#e74c3c}
.trial-msg.success{color:#4ade80}
.btn-disabled-trial{opacity:.45;cursor:not-allowed!important}
</style>
</head>
<body>

<div class="trial-bar" id="trial-bar">
<div class="trial-bar-msg">
<strong>Trial Required</strong>
You can view your existing habits, but creating, editing, or checking in is locked until you start a 14-day free trial.
</div>
<div class="trial-bar-actions">
<a class="btn-trial" href="https://stockyard.dev/" target="_blank" rel="noopener">Start 14-Day Trial</a>
<span class="trial-bar-divider">or</span>
<input type="text" class="key-input" id="trial-key-input" placeholder="SY-..." autocomplete="off" spellcheck="false">
<button class="btn-activate" id="trial-activate-btn" onclick="activateLicense()">Activate</button>
<span class="trial-msg" id="trial-msg"></span>
</div>
</div>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> TRAILHEAD</h1>
<button class="btn btn-p" onclick="openNew()">+ New Habit</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="progress" id="progress"></div>
<div class="toolbar">
<label class="arch-toggle"><input type="checkbox" id="arch-cb" onchange="toggleArchived()"> Show archived</label>
</div>
<div id="habitList"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='habits';

var customFields=[];
var habits=[],stats={},showArchived=false,editId=null,habitExtras={};

// ─── Helpers ──────────────────────────────────────────────────────

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var url=A+'/habits'+(showArchived?'?archived=true':'');
var resps=await Promise.all([
fetch(url).then(function(r){return r.json()}),
fetch(A+'/today').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
habits=resps[0].habits||[];
var todayView=resps[1]||{};
stats=resps[2]||{};

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
habitExtras=ex||{};
habits.forEach(function(h){
var x=habitExtras[h.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(h[k]===undefined)h[k]=x[k]});
});
}catch(e){habitExtras={}}

renderProgress(todayView);
renderStats();
}catch(e){
console.error('load failed',e);
habits=[];
}
renderHabits();
}

function renderProgress(today){
var done=today.done||0;
var total=today.total||0;
var pct=total?Math.round(done/total*100):0;
document.getElementById('progress').innerHTML=
'<div class="progress-text"><span>'+done+' of '+total+' completed today</span><span class="progress-pct">'+pct+'%</span></div>'+
'<div class="progress-bar"><div class="progress-fill" style="width:'+pct+'%"></div></div>'+
'<div class="progress-sub">'+esc(today.date||'')+'</div>';
}

function renderStats(){
var habits=stats.habits||0;
var totalChecks=stats.total_checks||0;
var activeStreaks=stats.active_streaks||0;
var rate=Math.round(stats.completion_rate||0);
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+habits+'</div><div class="st-l">Habits</div></div>'+
'<div class="st"><div class="st-v">'+totalChecks+'</div><div class="st-l">Total Checks</div></div>'+
'<div class="st"><div class="st-v">'+activeStreaks+'</div><div class="st-l">Active Streaks</div></div>'+
'<div class="st"><div class="st-v">'+rate+'%</div><div class="st-l">Today Rate</div></div>';
}

function renderHabits(){
var el=document.getElementById('habitList');
if(!habits.length){
var msg=window._emptyMsg||'No habits yet. Add one to start tracking.';
el.innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}
el.innerHTML=habits.map(habitHTML).join('');
}

function habitHTML(h){
var done=h.checked_today;
var cls='habit'+(h.archived?' archived':'');
var html='<div class="'+cls+'">';
if(window._trialRequired){
html+='<div class="check-btn '+(done?'done':'')+'" onclick="showTrialNudge()" title="Locked: trial required">'+(done?'&#10003;':'')+'</div>';
}else{
html+='<div class="check-btn '+(done?'done':'')+'" onclick="toggle(\''+h.id+'\','+done+')">'+(done?'&#10003;':'')+'</div>';
}
html+='<div class="color-dot" style="background:'+esc(h.color||'#c45d2c')+'"></div>';
html+='<div class="habit-info">';
html+='<div class="habit-name">'+esc(h.name);
if(h.archived)html+=' <span style="font-size:.5rem;color:var(--cm)">[archived]</span>';
html+='</div>';
if(h.description)html+='<div class="habit-desc">'+esc(h.description)+'</div>';
html+='<div class="habit-meta">';
html+='<span class="streak '+(h.streak>0?'active':'')+'">'+h.streak+'d streak</span>';
html+='<span>Best: '+h.best_streak+'d</span>';
html+='<span>'+h.total_checks+' total</span>';
html+='<span>'+esc(h.frequency||'daily')+'</span>';
html+='</div>';

// Custom field display
var customRows='';
customFields.forEach(function(f){
var v=h[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="habit-extra-row">';
customRows+='<span class="habit-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="habit-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)html+='<div class="habit-extra">'+customRows+'</div>';

html+='</div>';
if(!window._trialRequired){
html+='<div class="habit-actions">';
html+='<button class="icon-btn" onclick="openEdit(\''+h.id+'\')">edit</button>';
html+='</div>';
}
html+='</div>';
return html;
}

async function toggle(id,done){
try{
if(done){
await fetch(A+'/habits/'+id+'/uncheck',{method:'POST',headers:{'Content-Type':'application/json'},body:'{}'});
}else{
await fetch(A+'/habits/'+id+'/check',{method:'POST',headers:{'Content-Type':'application/json'},body:'{}'});
}
}catch(e){alert('Failed');return}
load();
}

function toggleArchived(){
showArchived=document.getElementById('arch-cb').checked;
load();
}

// ─── Modal: new / edit habit ──────────────────────────────────────

function renderExtrasInForm(values){
if(!customFields.length)return '';
var label=window._customSectionLabel||'Additional Details';
var h='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){
var v=values&&values[f.name]!==undefined?values[f.name]:'';
h+='<div class="fr"><label>'+esc(f.label)+'</label>';
if(f.type==='select'){
h+='<select id="ex-'+f.name+'">';
h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(String(o))+'</option>';
});
h+='</select>';
}else if(f.type==='textarea'){
h+='<textarea id="ex-'+f.name+'" rows="2">'+esc(String(v))+'</textarea>';
}else if(f.type==='number'||f.type==='integer'){
h+='<input type="number" id="ex-'+f.name+'" value="'+esc(String(v))+'">';
}else{
h+='<input type="text" id="ex-'+f.name+'" value="'+esc(String(v))+'">';
}
h+='</div>';
});
h+='</div>';
return h;
}

function modalHTML(habit){
var h=habit||{name:'',description:'',frequency:'daily',color:'#c45d2c',archived:false};
var isEdit=!!habit;
var nameph=window._placeholderName||'Exercise';
var html='<h2>'+(isEdit?'EDIT HABIT':'NEW HABIT')+'</h2>';
html+='<div class="fr"><label>Name *</label><input type="text" id="f-name" value="'+esc(h.name)+'" placeholder="'+esc(nameph)+'"></div>';
html+='<div class="fr"><label>Description</label><input type="text" id="f-desc" value="'+esc(h.description||'')+'" placeholder="optional"></div>';
html+='<div class="row2">';
html+='<div class="fr"><label>Frequency</label><select id="f-freq"><option value="daily"'+(h.frequency==='daily'?' selected':'')+'>Daily</option><option value="weekly"'+(h.frequency==='weekly'?' selected':'')+'>Weekly</option></select></div>';
html+='<div class="fr"><label>Color</label><input type="color" id="f-color" value="'+esc(h.color||'#c45d2c')+'"></div>';
html+='</div>';
if(isEdit){
html+='<div class="fr"><label><input type="checkbox" id="f-archived" '+(h.archived?'checked':'')+' style="width:auto;margin-right:.4rem">Archived</label></div>';
}
html+=renderExtrasInForm(h);
html+='<div class="acts">';
if(isEdit){
html+='<button class="btn btn-del" onclick="delHabit()">Delete</button>';
}
html+='<button class="btn" onclick="closeModal()">Cancel</button>';
html+='<button class="btn btn-p" onclick="save()">'+(isEdit?'Save':'Create')+'</button>';
html+='</div>';
return html;
}

function openNew(){
editId=null;
document.getElementById('mdl').innerHTML=modalHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-name');
if(n)n.focus();
}

function openEdit(id){
var h=null;
for(var i=0;i<habits.length;i++){if(habits[i].id===id){h=habits[i];break}}
if(!h)return;
editId=id;
document.getElementById('mdl').innerHTML=modalHTML(h);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function save(){
var nameEl=document.getElementById('f-name');
if(!nameEl||!nameEl.value.trim()){alert('Name is required');return}
var body={
name:nameEl.value.trim(),
description:document.getElementById('f-desc').value.trim(),
frequency:document.getElementById('f-freq').value,
color:document.getElementById('f-color').value
};
if(editId){
var arch=document.getElementById('f-archived');
body.archived=arch?arch.checked:false;
}

var extras={};
customFields.forEach(function(f){
var el=document.getElementById('ex-'+f.name);
if(!el)return;
var val;
if(f.type==='number'||f.type==='integer')val=parseFloat(el.value)||0;
else val=el.value.trim();
extras[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/habits/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/habits',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Create failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}
closeModal();
load();
}

async function delHabit(){
if(!editId)return;
if(!confirm('Delete this habit and all its check-ins?'))return;
await fetch(A+'/habits/'+editId,{method:'DELETE'});
closeModal();
load();
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.placeholder_name)window._placeholderName=cfg.placeholder_name;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
customFields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[]
});
});
}
}).catch(function(){
}).finally(function(){
checkTrialState();
load();
});
})();

// ─── trial-required license gating ───
window._trialRequired=false;

async function checkTrialState(){
try{
var resp=await fetch('/api/tier');
if(!resp.ok)return;
var data=await resp.json();
window._trialRequired=!!data.trial_required;
if(window._trialRequired){
document.getElementById('trial-bar').classList.add('show');
disableWriteControls();
// Re-render so check-btn/edit-btn handlers pick up trial state
if(typeof render==='function')render();
else if(typeof load==='function')load();
}else{
document.getElementById('trial-bar').classList.remove('show');
}
}catch(e){}
}

function disableWriteControls(){
var buttons=document.querySelectorAll('.hdr .btn, .hdr .btn-p');
buttons.forEach(function(b){
var t=b.textContent||'';
if(t.indexOf('New')!==-1||t.indexOf('Add')!==-1||t.indexOf('Habit')!==-1){
b.classList.add('btn-disabled-trial');
b.title='Locked: trial required';
b.onclick=function(e){
e.preventDefault();
showTrialNudge();
return false;
};
}
});
}

function showTrialNudge(){
var input=document.getElementById('trial-key-input');
if(input){
input.focus();
input.style.borderColor='var(--rust)';
setTimeout(function(){if(input)input.style.borderColor=''},1500);
}
}

async function activateLicense(){
var input=document.getElementById('trial-key-input');
var btn=document.getElementById('trial-activate-btn');
var msg=document.getElementById('trial-msg');
if(!input||!btn||!msg)return;
var key=(input.value||'').trim();
if(!key){
msg.className='trial-msg error';
msg.textContent='Paste your license key first';
input.focus();
return;
}
btn.disabled=true;
msg.className='trial-msg';
msg.textContent='Activating...';
try{
var resp=await fetch('/api/license/activate',{
method:'POST',
headers:{'Content-Type':'application/json'},
body:JSON.stringify({license_key:key})
});
var data=await resp.json();
if(!resp.ok){
msg.className='trial-msg error';
msg.textContent=data.error||'Activation failed';
btn.disabled=false;
return;
}
msg.className='trial-msg success';
msg.textContent='Activated. Reloading...';
setTimeout(function(){location.reload()},800);
}catch(e){
msg.className='trial-msg error';
msg.textContent='Network error: '+e.message;
btn.disabled=false;
}
}

document.addEventListener('DOMContentLoaded',function(){
var input=document.getElementById('trial-key-input');
if(input){
input.addEventListener('keydown',function(e){
if(e.key==='Enter')activateLicense();
});
}
});
</script>
</body>
</html>`
