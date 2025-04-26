'use strict';

const express = require('express');
const app = express();
let path = require('path');
let sdk = require('./sdk');

const PORT = 8001;
const HOST = '0.0.0.0';

// body parsing
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// 📌 static 파일 서비스 (index.html, app.js 파일 제공)
app.use(express.static(path.join(__dirname, '../client')));

// 📌 루트("/") 접속시 index.html 보내기
app.get('/', function(req, res){
    res.sendFile(path.join(__dirname, '../client/index.html'));
});

// ================= 기존 ABStore API ==================

// Init
app.get('/init', function (req, res) {
    let a = req.query.a;
    let aval = req.query.aval;
    let b = req.query.b;
    let bval = req.query.bval;
    let args = [a, aval, b, bval];
    sdk.send(false, 'Init', args, res);
});

// Invoke
app.get('/invoke', function (req, res) {
    let a = req.query.a;
    let b = req.query.b;
    let value = req.query.value;
    let args = [a, b, value];
    sdk.send(false, 'Invoke', args, res);
});

// Query
app.get('/query', function (req, res) {
    let name = req.query.name;
    let args = [name];
    sdk.send(true, 'Query', args, res);
});

// Query All
app.get('/queryAll', function (req, res) {
    sdk.send(true, 'GetAllQuery', [], res);
});

// Delete
app.get('/delete', function (req, res) {
    let name = req.query.name;
    let args = [name];
    sdk.send(false, 'Delete', args, res);
});

// ================= 깐부 대출 Loan API ==================

// Create Loan Request
app.get('/createLoan', function (req, res) {
    let { id, requester, amount, durationDays } = req.query;
    let args = [id, requester, amount, durationDays];
    sdk.send(false, 'CreateLoanRequest', args, res);
});

// Approve Loan Request
app.get('/approveLoan', function (req, res) {
    let { id, provider } = req.query;
    let args = [id, provider];
    sdk.send(false, 'ApproveLoanRequest', args, res);
});

// Delete Loan Request
app.get('/deleteLoan', function (req, res) {
    let { id } = req.query;
    let args = [id];
    sdk.send(false, 'DeleteLoanRequest', args, res);
});

// Update Loan Request
app.get('/updateLoan', function (req, res) {
    let { id, newAmount, newDurationDays } = req.query;
    let args = [id, newAmount, newDurationDays];
    sdk.send(false, 'UpdateLoanRequest', args, res);
});

// Query Loan Request
app.get('/queryLoan', function (req, res) {
    let { id } = req.query;
    let args = [id];
    sdk.send(true, 'QueryLoanRequest', args, res);
});

// Query All Loan Requests
app.get('/queryAllLoans', function (req, res) {
    sdk.send(true, 'QueryAllLoanRequests', [], res);
});

// ================= 서버 시작 ==================
app.listen(PORT, HOST);
console.log(`🌟 Server running at http://${HOST}:${PORT}/`);
