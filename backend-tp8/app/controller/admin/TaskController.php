<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\model\Task;
use app\service\TaskService;

class TaskController extends BaseController
{
    public function index()
    {
        $status = $this->request->get('status');
        $keyword = trim((string) $this->request->get('keyword', ''));
        $page = (int) $this->request->get('page', 1);
        $pageSize = (int) $this->request->get('page_size', 10);

        $query = Task::order('id', 'desc');
        if ($status !== null && $status !== '' && (int) $status >= 0) {
            $query->where('status', (int) $status);
        }
        if ($keyword !== '') {
            $query->whereLike('title', '%' . $keyword . '%');
        }

        $pager = $query->paginate(['list_rows' => $pageSize, 'page' => $page])->toArray();
        return json_success([
            'list' => $pager['data'] ?? [],
            'total' => (int) ($pager['total'] ?? 0),
            'page' => (int) ($pager['current_page'] ?? $page),
            'page_size' => (int) ($pager['per_page'] ?? $pageSize),
        ]);
    }

    public function create()
    {
        $data = [
            'title' => trim((string) $this->request->post('title', '')),
            'summary' => trim((string) $this->request->post('summary', '')),
            'detail' => (string) $this->request->post('detail', ''),
            'doc_url' => (string) $this->request->post('doc_url', ''),
            'reward_points' => (int) $this->request->post('reward_points', 0),
            'status' => TaskService::TASK_STATUS_ONLINE,
        ];
        try {
            validate(\app\validate\TaskCreateValidate::class)->check($data);
        } catch (\think\exception\ValidateException $e) {
            return json_error($e->getError());
        }

        $task = Task::create($data);
        return json_success($task->toArray());
    }

    public function update(int $id)
    {
        $task = Task::where('id', $id)->find();
        if (!$task) {
            return json_error('任务不存在');
        }

        $task->title = trim((string) $this->request->post('title', $task->title));
        $task->summary = trim((string) $this->request->post('summary', $task->summary));
        $task->detail = (string) $this->request->post('detail', $task->detail);
        $task->doc_url = (string) $this->request->post('doc_url', $task->doc_url);
        $task->reward_points = (int) $this->request->post('reward_points', $task->reward_points);
        $task->status = (int) $this->request->post('status', $task->status);
        $task->save();

        return json_success($task->toArray());
    }

    public function off(int $id)
    {
        $task = Task::where('id', $id)->find();
        if (!$task) {
            return json_error('任务不存在');
        }

        $task->status = TaskService::TASK_STATUS_OFFLINE;
        $task->save();

        return json_success();
    }
}
