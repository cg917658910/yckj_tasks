<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\service\WithdrawalService;
use think\facade\Db;

class WithdrawalController extends BaseController
{
    public function index()
    {
        $status = $this->request->get('status');
        $keyword = trim((string) $this->request->get('keyword', ''));
        $page = (int) $this->request->get('page', 1);
        $pageSize = (int) $this->request->get('page_size', 10);

        $query = Db::name('withdrawals')->alias('w')
            ->leftJoin('users u', 'w.user_id = u.id')
            ->field('w.*, u.username');

        if ($status !== null && $status !== '' && (int) $status >= 0) {
            $query->where('w.status', (int) $status);
        }
        if ($keyword !== '') {
            $query->whereLike('u.username', '%' . $keyword . '%');
        }

        $pager = $query->order('w.id', 'desc')->paginate(['list_rows' => $pageSize, 'page' => $page])->toArray();
        return json_success([
            'list' => $pager['data'] ?? [],
            'total' => (int) ($pager['total'] ?? 0),
            'page' => (int) ($pager['current_page'] ?? $page),
            'page_size' => (int) ($pager['per_page'] ?? $pageSize),
        ]);
    }

    public function pay(int $id)
    {
        $adminId = (int) ($this->request->admin['admin_id'] ?? 0);

        try {
            WithdrawalService::approve($adminId, $id);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success();
    }

    public function reject(int $id)
    {
        $reason = trim((string) $this->request->post('reason', ''));
        if ($reason === '') {
            return json_error('请填写驳回原因');
        }

        $adminId = (int) ($this->request->admin['admin_id'] ?? 0);

        try {
            WithdrawalService::reject($adminId, $id, $reason);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success();
    }
}
