typedef enum {
    GN_AST_NUMBER,
    GN_AST_BOOLEAN,
    GN_AST_SYMBOL,

    GN_AST_ASSIGN,
    GN_AST_ADD,
    GN_AST_SUBTRACT,
    GN_AST_MULTIPLY,
    GN_AST_DIVIDE,
    GN_AST_COMPARE,
    GN_AST_INVERSE_COMPARE
} gn_ast_node_type_t;

typedef struct gn_ast_node_t {
    gn_ast_node_type_t node_type;
    int value;
    struct gn_ast_node_t **children;
    int num_children;
} gn_ast_node_t;

typedef struct {
    int error;

    gn_ast_node_t *stack[1024];
    int position;

    char *symbols[1024];
    int next_symbol;

    char *input_buffer;
    int buffer_len;
    int buffer_offset;
} gn_parser_context_t;

gn_parser_context_t *gn_global_context();
gn_ast_node_t *gn_parse(gn_parser_context_t *context, char *buffer, int len);
char *gn_get_symbol(gn_parser_context_t *context, gn_ast_node_t *node);