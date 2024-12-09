module alu #(
    parameter N = 8
) (
    input logic [N-1:0] a, b,
    input logic [2:0] op,
    input logic carry_in,
    output logic [N-1:0] c,
    output logic a_gt_b, a_eq_b, c_eq_zero, carry_out
);

    logic [N-1:0] r;
    logic co;

    always_comb begin
        r = 0;
        co = 0;

        a_gt_b = (a > b);
        a_eq_b = (a == b);

        case (op)
            3'b000 : {co, r} = a + b + { 8'b0 , carry_in};
            3'b001 : {r, co} = {carry_in, a};
            3'b010 : {co, r} = {a, carry_in};
            3'b011 : r = ~a;
            3'b100 : r = a & b;
            3'b101 : r = a | b;
            3'b110 : r = a ^ b;
            default : r = 0;
        endcase

        c_eq_zero = (r == 0);
        c = r;
        carry_out = co;
    end

endmodule