<?php

declare(strict_types=1);

use think\facade\Route;
use app\middleware\AdminAuth;
use app\middleware\UserAuth;

// Admin auth (public)
Route::post('admin/auth/login', 'admin.AuthController/login');
Route::post('admin/auth/logout', 'admin.AuthController/logout')->middleware(AdminAuth::class);

// Admin protected routes
Route::group('admin', function () {
    Route::get('points/rules', 'admin.PointsController/rules');
    Route::put('points/rules', 'admin.PointsController/updateRules');

    Route::get('tasks', 'admin.TaskController/index');
    Route::post('tasks', 'admin.TaskController/create');
    Route::put('tasks/:id', 'admin.TaskController/update');
    Route::put('tasks/:id/off', 'admin.TaskController/off');

    Route::get('claims', 'admin.ClaimController/index');
    Route::post('claims/:id/approve', 'admin.ClaimController/approve');
    Route::post('claims/:id/reject', 'admin.ClaimController/reject');

    Route::get('withdrawals', 'admin.WithdrawalController/index');
    Route::post('withdrawals/:id/pay', 'admin.WithdrawalController/pay');
    Route::post('withdrawals/:id/reject', 'admin.WithdrawalController/reject');

    Route::get('users', 'admin.UserController/index');
    Route::put('users/:id/status', 'admin.UserController/status');
    Route::post('users/:id/points', 'admin.UserController/adjustPoints');
    Route::get('users/:id/tasks', 'admin.UserController/tasks');
    Route::get('users/:id/withdrawals', 'admin.UserController/withdrawals');

    Route::post('upload/image', 'admin.UploadController/image');
})->middleware(AdminAuth::class);

// User auth (public)
Route::post('user/auth/register', 'user.AuthController/register');
Route::post('user/auth/login', 'user.AuthController/login');
Route::post('user/auth/logout', 'user.AuthController/logout')->middleware(UserAuth::class);

// User protected routes
Route::group('user', function () {
    Route::get('tasks', 'user.TaskController/index');
    Route::get('tasks/:id', 'user.TaskController/detail');
    Route::post('tasks/:id/claim', 'user.TaskController/claim');

    Route::get('claims/current', 'user.ClaimController/current');
    Route::get('claims/history', 'user.ClaimController/history');
    Route::post('claims/:id/submit', 'user.ClaimController/submit');

    Route::get('profile', 'user.ProfileController/info');
    Route::put('profile/password', 'user.ProfileController/changePassword');
    Route::put('profile/wechat-qr', 'user.ProfileController/updateWechatQr');

    Route::post('withdrawals', 'user.WithdrawalController/apply');
    Route::get('withdrawals', 'user.WithdrawalController/index');

    Route::post('upload/image', 'user.UploadController/image');
})->middleware(UserAuth::class);
