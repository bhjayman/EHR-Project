//SPDX-License-Identifier: Apache-2.0

var ehr = require('./EHRcontroller.js');

module.exports = function(app){

  app.get('/get_data/:id', function(req, res){
    ehr.get_data(req, res);
  });

  app.get('/req_access/:idMedecin/:idPatient', function(req, res){
    ehr.req_access(req, res);
  });

  app.get('/get_all_access/:id', function(req, res){
    ehr.get_all_access(req, res);
  });

  app.get('/get_all_pat_access/:idP/:idM/:type', function(req, res){
    ehr.get_all_pat_access(req, res);
  });

  app.get('/get_perm/:idM/:idP', function(req, res){
    ehr.get_perm(req, res);
  });

  app.get('/get_all', function(req, res){
    ehr.get_all(req, res);
  });

  app.get('/edit_perm/:idP/:idM/:dateAut/:statut', function(req, res){
    ehr.edit_perm(req, res);
  });

  app.get('/get_dossier/:idM/:idP/:dateAut/:mode', function(req, res){
    ehr.get_dossier(req, res);
  });

  app.get('/add_to_dossier/:id/:data', function(req, res){
    ehr.add_to_dossier(req, res);
  });

  app.get('/add_medecin/:id/:nomp/:sp', function(req, res){
    ehr.add_medecin(req, res);
  });

  app.get('/add_patient/:id/:nomp/:dn', function(req, res){
    ehr.add_patient(req, res);
  });

  app.get('/add_agent/:id/:nomp/:lab', function(req, res){
    ehr.add_agent(req, res);
  });

}
