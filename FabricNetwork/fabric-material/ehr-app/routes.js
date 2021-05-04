//SPDX-License-Identifier: Apache-2.0

var tuna = require('./controller.js');
var ehr = require('./EHRcontroller.js');

module.exports = function(app){

  app.get('/get_tuna/:id', function(req, res){
    tuna.get_tuna(req, res);
  });
  app.get('/add_tuna/:tuna', function(req, res){
    tuna.add_tuna(req, res);
  });
  app.get('/get_all_tuna', function(req, res){
    tuna.get_all_tuna(req, res);
  });
  app.get('/change_holder/:holder', function(req, res){
    tuna.change_holder(req, res);
  });

  app.get('/test',function(req,res){
    res.json({ username: 'Ayman' })
  })

////////////////////////////////////////////////////////////////

  app.get('/check_user/:id', function(req, res){
    ehr.check_user(req, res);
  });

  app.get('/req_access/:idMedecin/:idPatient', function(req, res){
    ehr.req_access(req, res);
  });

  app.get('/get_all_access/:id', function(req, res){
    ehr.get_all_access(req, res);
  });

  app.get('/get_perm/:idM/:idP', function(req, res){
    ehr.get_perm(req, res);
  });

  app.get('/get_all', function(req, res){
    ehr.get_all(req, res);
  });

}
