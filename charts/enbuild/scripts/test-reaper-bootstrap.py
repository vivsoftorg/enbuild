#!/usr/bin/env python3
"""Exercise the complete production bootstrap with a recording kubectl stub."""
import json, os, subprocess, sys, tempfile, unittest
from pathlib import Path
SCRIPT=Path(__file__).with_name('create-bootstrap-secrets.sh')
STUB='''#!PYTHON
import base64,json,os,sys
from pathlib import Path
a=sys.argv[1:]
with Path(os.environ['CALL_LOG']).open('a') as f:f.write(json.dumps(a)+'\\n')
if 'get' in a and 'namespace' in a:sys.exit(0)
if 'get' in a and 'secret' in a:
 if os.environ.get('EXISTING')=='1':
  if any('rabbitmq-password' in x for x in a):print(base64.b64encode(b'fixture-rabbit-password').decode(),end='')
  sys.exit(0)
 sys.exit(1)
if 'create' in a:print('apiVersion: v1\\nkind: Secret\\nmetadata:\\n  name: fixture')
if 'apply' in a:sys.stdin.read()
'''.replace('PYTHON',sys.executable)
class ReaperBootstrap(unittest.TestCase):
 def run_script(self,values=None,existing=False):
  with tempfile.TemporaryDirectory()as tmp:
   p=Path(tmp);stub=p/'kubectl';stub.write_text(STUB);stub.chmod(0o755);log=p/'calls.jsonl'
   env={k:v for k,v in os.environ.items()if not k.startswith(('REAPER_','REPO1_','GITLAB_','IA_','OBS_'))and k not in['FORCE','NAMESPACE','RELEASE']}
   env.update(PATH=str(p)+os.pathsep+env['PATH'],CALL_LOG=str(log),TMPDIR=tmp,EXISTING='1'if existing else'0')
   env.update(values or {});r=subprocess.run(['bash',str(SCRIPT)],env=env,capture_output=True,text=True,timeout=20)
   return r,[json.loads(x)for x in log.read_text().splitlines()]if log.exists()else[]
 def valid(self):return {'REAPER_ACCESS_KEY_ID':'fixture-access-only','REAPER_SECRET_ACCESS_KEY':'fixture-secret-only','REAPER_ACCOUNT_ID':'123456789012'}
 def test_incomplete_or_malformed_refuses_before_all_kubernetes_calls(self):
  cases=[{'REAPER_ACCESS_KEY_ID':'fixture-access-only'},{'REAPER_SECRET_ACCESS_KEY':'fixture-secret-only'},{'REAPER_ACCOUNT_ID':'123456789012'},{'REAPER_RESOURCE_TAG_KEY':'vendor'},dict(self.valid(),REAPER_SECRET_ACCESS_KEY='')]
  cases +=[dict(self.valid(),REAPER_ACCOUNT_ID=x)for x in['','12345678901','1234567890123','12345678901x',' 123456789012','123456789012\n']]
  for v in cases:
   with self.subTest(v=list(v)):
    r,calls=self.run_script(v);self.assertNotEqual(r.returncode,0);self.assertEqual(calls,[]);self.assertNotIn('fixture-access-only',r.stdout+r.stderr);self.assertNotIn('fixture-secret-only',r.stdout+r.stderr)
 def test_partial_resource_tag_pair_refuses_before_calls(self):
  for k in['REAPER_RESOURCE_TAG_KEY','REAPER_RESOURCE_TAG_VALUE']:
   r,calls=self.run_script(dict(self.valid(),**{k:'fixture'}));self.assertNotEqual(r.returncode,0);self.assertEqual(calls,[])
 def test_omission_preserves_optional_reaper(self):
  r,calls=self.run_script();self.assertEqual(r.returncode,0,r.stderr);self.assertFalse(any('enbuild-ib-reaper-svc'in x for x in calls));self.assertTrue(any('apply'in x for x in calls))
 def test_complete_exact_account_is_projected(self):
  r,calls=self.run_script(self.valid());self.assertEqual(r.returncode,0,r.stderr)
  creates=[x for x in calls if 'create'in x and'enbuild-ib-reaper-svc'in x];self.assertEqual(len(creates),1)
  literals={x.split('=',2)[1]:x.split('=',2)[2]for x in creates[0]if x.startswith('--from-literal=')}
  self.assertEqual(literals,{'REAPER_SVC_AWS_ACCESS_KEY_ID':'fixture-access-only','REAPER_SVC_AWS_SECRET_ACCESS_KEY':'fixture-secret-only','REAPER_SVC_AWS_ACCOUNT_ID':'123456789012'})
  self.assertNotIn('fixture-secret-only',r.stdout+r.stderr)
 def test_complete_custom_tag_pair_is_preserved(self):
  r,calls=self.run_script(dict(self.valid(),REAPER_RESOURCE_TAG_KEY='fixture-approved',REAPER_RESOURCE_TAG_VALUE='yes'));self.assertEqual(r.returncode,0,r.stderr)
  c=next(x for x in calls if 'create'in x and'enbuild-ib-reaper-svc'in x)
  self.assertIn('--from-literal=REAPER_SVC_RESOURCE_TAG_KEY=fixture-approved',c);self.assertIn('--from-literal=REAPER_SVC_RESOURCE_TAG_VALUE=yes',c)
 def test_existing_reaper_secret_kept(self):
  r,calls=self.run_script(self.valid(),existing=True);self.assertEqual(r.returncode,0,r.stderr);self.assertTrue(any('get'in x and'enbuild-ib-reaper-svc'in x for x in calls));self.assertFalse(any('create'in x and'enbuild-ib-reaper-svc'in x for x in calls))
if __name__=='__main__':unittest.main()
