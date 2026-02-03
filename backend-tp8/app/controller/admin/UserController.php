<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\service\PointsService;
use think\facade\Db;

class UserController extends BaseController
{
    public function index()
    {
        $status = $this->request->get('status');
        $keyword = trim((string) $this->request->get('keyword', ''));
        $page = (int) $this->request->get('page', 1);
        $pageSize = (int) $this->request->get('page_size', 10);

        $query = Db::name('users')->alias('u')
            ->leftJoin('points_accounts p', 'u.id = p.user_id')
            ->leftJoin('user_profiles up', 'u.id = up.user_id')
            ->field('u.id,u.username,u.status,u.created_at,p.available_points,p.frozen_points,p.withdrawn_points,up.total_withdrawn')
            ->order('u.id', 'desc');

        if ($status !== null && $status !== '' && (int) $status >= 0) {
            $query->where('u.status', (int) $status);
        }
        if ($keyword !== '') {
            $query->whereLike('u.username', '%' . $keyword . '%');
        }

        $pager = $query->paginate(['list_rows' => $pageSize, 'page' => $page])->toArray();

        return json_success([
            'list' => $pager['data'] ?? [],
            'total' => (int) ($pager['total'] ?? 0),
            'page' => (int) ($pager['current_page'] ?? $page),
            'page_size' => (int) ($pager['per_page'] ?? $pageSize),
        ]);
    }

    public function status(int $id)
    {
        $data = $this->request->only(['status'], 'post');
        try {
            validate(\app\validate\UserStatusValidate::class)->check($data);
        } catch (\think\exception\ValidateException $e) {
            return json_error($e->getError());
        }

        $status = (int) $data['status'];
        Db::name('users')->where('id', $id)->update(['status' => $status]);
        return json_success();
    }

    public function adjustPoints(int $id)
    {
        $points = (int) $this->request->post('points', 0);
        $remark = trim((string) $this->request->post('remark', '手动调整积分'));

        if ($points === 0) {
            return json_error('积分调整值不能为 0');
        }

        try {
            PointsService::addAvailable($id, $points, 'manual_adjust', null, $remark);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success();
    }

    public function tasks(int $id)
    {
        $list = Db::name('task_claims')->alias('c')
            ->leftJoin('tasks t', 'c.task_id = t.id')
            ->field('c.*, t.title as task_title, t.reward_points')
            ->where('c.user_id', $id)
            ->order('c.id', 'desc')
            ->select()
            ->toArray();

        return json_success(['list' => $list]);
    }

    public function withdrawals(int $id)
    {
        $list = Db::name('withdrawals')->where('user_id', $id)->order('id', 'desc')->select()->toArray();
        return json_success(['list' => $list]);
    }
}
